package template

const (
	Help  = `help placeholder for {{folderName}}`
	Umask = `#!/bin/sh
umask 0011
$1
`
)
