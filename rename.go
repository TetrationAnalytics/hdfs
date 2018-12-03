package hdfs

import (
	"errors"
	"os"

	hdfs "github.com/colinmarc/hdfs/protocol/hadoop_hdfs"
	"github.com/colinmarc/hdfs/rpc"
	"github.com/golang/protobuf/proto"
)

// Rename renames (moves) a file.
func (c *Client) Rename(oldpath, newpath string) error {
	return c.RenameWithOverwriteOption(oldpath, newpath, true)
}

// RenameWithOverwrite renames (moves) a file. Overwrite option is taken as input.
func (c *Client) RenameWithOverwriteOption(oldpath, newpath string, overwrite bool) error {
	f, err := c.getFileInfo(newpath)
	if err != nil && !os.IsNotExist(err) {
		return &os.PathError{"rename", newpath, err}
	}

	// If overwrite is not enabled and destPath exists throw error.
	if !overwrite && f != nil {
		return &os.PathError{"rename", newpath, errors.New("Path exists.")}
	}

	req := &hdfs.Rename2RequestProto{
		Src:           proto.String(oldpath),
		Dst:           proto.String(newpath),
		OverwriteDest: proto.Bool(overwrite),
	}
	resp := &hdfs.Rename2ResponseProto{}

	err = c.namenode.Execute("rename2", req, resp)
	if err != nil {
		if nnErr, ok := err.(*rpc.NamenodeError); ok {
			err = interpretException(nnErr.Exception, err)
		}

		return &os.PathError{"rename", oldpath, err}
	}

	return nil
}
