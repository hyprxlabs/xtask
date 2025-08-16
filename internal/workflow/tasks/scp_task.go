package tasks

import (
	"context"
	"io"
	"net"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/hyprxlabs/xtasks/internal/errors"
	"github.com/hyprxlabs/xtasks/internal/schema"
	goph "github.com/melbahja/goph"
	"golang.org/x/crypto/ssh"
)

func runSCP(ctx TaskContext) *TaskResult {

	//https://github.com/melbahja/goph

	res := NewTaskResult()

	if len(ctx.Task.Files) == 0 {
		return res.Fail(errors.New("No files specified for SCP task"))
	}

	uri, err := url.Parse(ctx.Task.Uses)
	if err != nil {
		return res.Fail(errors.New("Invalid SSH URI: " + err.Error()))
	}

	if uri.Scheme != "scp" {
		return res.Fail(errors.New("Invalid SSH URI scheme: " + uri.Scheme))
	}

	direction := uri.Path
	if direction == "" {
		direction = "upload"
	}

	targets := []schema.SshTarget{}
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
			password = ctx.Task.Env[password]
		}

		identity := uri.Query().Get("identity")

		targets = append(targets, schema.SshTarget{
			Host:     uri.Host,
			User:     &user,
			Port:     &port,
			Identity: &identity,
			Password: &password,
		})
	} else if len(ctx.Task.Targets) > 0 {
		targetNames := ctx.Task.Targets

		for _, targetName := range targetNames {
			target, ok := ctx.Targets[targetName]
			if ok {
				targets = append(targets, target)
			} else {
				for _, target := range ctx.Targets {
					for _, group := range target.Groups {
						if group == targetName {
							targets = append(targets, target)
						}
					}
				}
			}
		}
	} else {
		for _, value := range ctx.Targets {
			targets = append(targets, value)
		}
	}

	if len(targets) == 0 {
		return res.Fail(errors.New("No targets found for SSH task"))
	}

	for _, target := range targets {
		if err := runScpTarget(ctx.Context, direction, ctx, target); err != nil {
			return res.Fail(err)
		}
	}

	res.End()

	// Placeholder for running an SSH command
	// This would typically involve executing the command over SSH
	return res.Ok()
}

func runScpTarget(ctx context.Context, direction string, taskContext TaskContext, target schema.SshTarget) error {
	var auth goph.Auth
	var err error
	identity := ""
	password := ""
	if target.Identity != nil {
		identity = *target.Identity
	}
	if target.Password != nil {
		password = *target.Password
		if password != "" {
			password = taskContext.Task.Env[password]
		}
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
		return errors.New("Failed to create SSH authentication: " + err.Error())
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

	for _, file := range taskContext.Task.Files {
		parts := strings.Split(file, ":")
		if len(parts) != 2 {
			err2 := errors.New("Invalid SCP file format, expected 'source:destination'")
			err2 = errors.WithCause(err2, err)
			return err2
		}
		source := parts[0]
		destination := parts[1]

		if direction == "upload" {
			err = Upload(ctx, client, source, destination)
		}
		if direction == "download" {
			err = Download(ctx, client, destination, source)
		}

		if err != nil {
			err2 := errors.New("Failed to transfer file " + source + " to " + destination + ": " + err.Error())
			err2 = errors.WithCause(err2, err)
			return err2
		}
	}

	return nil
}

type readerFunc func(p []byte) (n int, err error)

func (rf readerFunc) Read(p []byte) (n int, err error) { return rf(p) }

func Upload(ctx context.Context, c *goph.Client, localPath string, remotePath string) (err error) {

	local, err := os.Open(localPath)
	if err != nil {
		return
	}
	defer local.Close()

	ftp, err := c.NewSftp()
	if err != nil {
		return
	}
	defer ftp.Close()

	remote, err := ftp.Create(remotePath)
	if err != nil {
		return
	}
	defer remote.Close()

	_, err = io.Copy(remote, readerFunc(func(p []byte) (int, error) {

		// golang non-blocking channel: https://gobyexample.com/non-blocking-channel-operations
		select {

		// if context has been canceled
		case <-ctx.Done():
			// stop process and propagate "context canceled" error
			return 0, ctx.Err()
		default:
			// otherwise just run default io.Reader implementation
			return local.Read(p)
		}
	}))
	return
}

// Download file from remote server!
func Download(ctx context.Context, c *goph.Client, remotePath string, localPath string) (err error) {

	local, err := os.Create(localPath)
	if err != nil {
		return
	}
	defer local.Close()

	ftp, err := c.NewSftp()
	if err != nil {
		return
	}
	defer ftp.Close()

	remote, err := ftp.Open(remotePath)
	if err != nil {
		return
	}
	defer remote.Close()

	_, err = io.Copy(local, readerFunc(func(p []byte) (int, error) {

		// golang non-blocking channel: https://gobyexample.com/non-blocking-channel-operations
		select {

		// if context has been canceled
		case <-ctx.Done():
			// stop process and propagate "context canceled" error
			return 0, ctx.Err()
		default:
			// otherwise just run default io.Reader implementation
			return remote.Read(p)
		}
	}))
	if err != nil {
		return err
	}

	return local.Sync()
}
