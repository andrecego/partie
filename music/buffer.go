package music

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/jonas747/dca"
)

func addToBuffer(song Song) error {
	fmt.Println("Adding to buffer: ", song.GetPath())

	file, err := dcaEncode(song.GetPath())
	if err != nil {
		fmt.Println("Error opening opus file :", err)
		return err
	}

	var opuslen int16

	for {
		// Read opus frame length from dca file.
		err = binary.Read(file, binary.LittleEndian, &opuslen)

		// If this is the end of the file, just return.
		// if err == io.EOF || err == io.ErrUnexpectedEOF {
		if err == io.EOF {
			return nil
		}

		if err != nil {
			fmt.Println("Error reading from dca file :", err)
			return err
		}

		// Read encoded pcm from dca file.
		InBuf := make([]byte, opuslen)
		err = binary.Read(file, binary.LittleEndian, &InBuf)

		// Should not be any end of file errors
		if err != nil {
			fmt.Println("Error reading from dca file :", err)
			return err
		}

		// Append encoded pcm data to the buffer.
		currentDJ.Buffer = append(currentDJ.Buffer, InBuf)
	}
}

func dcaEncode(fileName string) (*dca.EncodeSession, error) {
	opts := dca.StdEncodeOptions
	opts.RawOutput = true
	opts.Bitrate = 96
	opts.Application = "lowdelay"

	return dca.EncodeFile(fileName, opts)
}
