package tasks

import (
	"context"
	"net"
	"net/url"
	"os"
	"strconv"

	"github.com/hyprxlabs/xtask/errors"
	"github.com/hyprxlabs/xtask/types"
	goph "github.com/melbahja/goph"
	"golang.org/x/crypto/ssh"
)

func runSSH(ctx TaskContext) *TaskResult {
	//https://github.com/melbahja/goph

	res := NewTaskResult()
	uses := ctx.Data.Uses
	if uses == "ssh" {
		uses = "ssh://"
	}

	uri, err := url.Parse(uses)
	if err != nil {
		return res.Fail(errors.New("Invalid SSH URI: " + err.Error()))
	}

	if uri.Scheme != "ssh" {
		return res.Fail(errors.New("Invalid SSH URI scheme: " + uri.Scheme))
	}

	targets := []types.Host{}
	if uri.Host != "" {
		user := ""
		if uri.User != nil {
			user = uri.User.Username()
		}

		port := 22
		if uri.Port() != "" {
			port, err = strconv.Atoi(uri.Port())
			if err != nil {
				return res.Fail(errors.New("Invalid port in SSH URI: " + err.Error()))
			}
		}

		password, ok := uri.User.Password()
		if ok && password != "" {
			p, ok := ctx.Data.Env.Get(password)
			if ok {
				password = p
			}
		}

		identity := uri.Query().Get("identity")
		targets = append(targets, types.Host{
			Host:     uri.Host,
			User:     &user,
			Port:     &port,
			Identity: &identity,
			Password: &password,
		})
	} else {
		for _, t := range ctx.Data.Hosts {
			targets = append(targets, types.Host{
				Host:     t.Host,
				User:     t.User,
				Port:     t.Port,
				Identity: t.Identity,
				Password: t.Password,
			})
		}
	}

	if len(targets) == 0 {
		return res.Fail(errors.New("No targets found for SSH task"))
	}

	for _, target := range targets {
		err := runSSHTarget(ctx.Context, ctx, target)
		if errors.Is(err, context.Canceled) {
			return res.Cancel("Task " + ctx.Task.Id + " cancelled")
		}

		if errors.Is(err, context.DeadlineExceeded) {
			return res.Cancel("Task " + ctx.Task.Id + " cancelled due to timeout")
		}
	}

	res.End()

	// Placeholder for running an SSH command
	// This would typically involve executing the command over SSH
	return res.Ok()
}

type SshRun struct {
	Error error
}

func runSSHTarget(ctx context.Context, taskContext TaskContext, target types.Host) error {
	signal := make(chan SshRun)

	var auth goph.Auth
	var err error
	identity := ""
	password := ""
	run := ""
	if target.Identity != nil {
		identity = *target.Identity
	}

	if target.Password != nil {
		password = *target.Password
	}

	if identity == "" && password != "" {
		auth = goph.Password(password)
	} else if goph.HasAgent() {
		auth, err = goph.UseAgent()
	} else if identity != "" {
		auth, err = goph.Key(identity, password)
	} else {
		return errors.New("No authentication method provided for SSH task")
	}

	if err != nil {
		signal <- SshRun{Error: errors.New("Failed to create SSH authentication: " + err.Error())}
	}

	port := 22
	if target.Port != nil && *target.Port > 0 {
		port = int(*target.Port)
	}
	user := ""
	if target.User != nil && *target.User != "" {
		user = *target.User
	}

	client, err := goph.NewConn(&goph.Config{
		User: user,
		Addr: target.Host,
		Port: uint(port),
		Auth: auth,
		Callback: func(host string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	})

	if err != nil {
		err2 := errors.New("Failed to connect to SSH target " + target.Host + ": " + err.Error())
		err2 = errors.WithCause(err2, err)
		return err2
	}

	defer client.Close()

	var sess *ssh.Session

	if sess, err = client.NewSession(); err != nil {
		err2 := errors.New("Failed to create SSH session: " + err.Error())
		return err2
	}

	defer sess.Close()

	go func() {

		if taskContext.Data.Env.Len() < 0 {
			// only set env values that are explicitly set in the task
			for _, key := range taskContext.Task.Env.Keys() {
				value, _ := taskContext.Task.Env.Get(key)
				sess.Setenv(key, value)
			}
		}

		sess.Stdout = os.Stdout
		sess.Stderr = os.Stderr
		err = sess.Run(run)

		if err != nil {
			err2 := errors.New("Failed to run command on SSH target " + target.Host + ": " + err.Error())
			err2 = errors.WithCause(err2, err)
			signal <- SshRun{Error: err2}
			return
		}

		signal <- SshRun{Error: nil}
	}()

	select {
	case <-ctx.Done():
		sess.Signal(ssh.SIGINT)
		return ctx.Err()
	case result := <-signal:
		return result.Error
	}
}
