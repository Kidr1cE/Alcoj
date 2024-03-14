package master

import (
	"context"
	"io"
	"log"
	"os"

	pb "alcoj/proto"
)

func uploadFile(client pb.SandboxClient, filename string, path string) error {
	stream, err := client.SendRequirements(context.Background())
	if err != nil {
		return err
	}

	file, err := os.Open(path)
	if err != nil {
		log.Println(err)
		return err
	}
	defer file.Close()

	buf := make([]byte, 1024*1024)
	for {
		n, err := file.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		if err := stream.Send(&pb.FileChunk{
			Filename:    filename,
			Content:     buf[:n],
			IsLastChunk: false,
		}); err != nil {
			return err
		}
		resp, err := stream.Recv()
		if err != nil {
			return err
		}
		log.Println("recv: ", resp)
	}

	stream.Send(&pb.FileChunk{
		Filename:    filename,
		Content:     nil,
		IsLastChunk: true,
	})
	stream.Recv()

	return nil
}
