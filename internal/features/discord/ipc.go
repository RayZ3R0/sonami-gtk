package discord

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"sync/atomic"
)

const (
	opcodeHandshake uint32 = 0
	opcodeFrame     uint32 = 1
	opcodeClose     uint32 = 2
)

var nonceSeq atomic.Int64

func nextNonce() string {
	return fmt.Sprintf("%d", nonceSeq.Add(1))
}

type ipcClient struct {
	conn net.Conn
}

func dialIPC() (*ipcClient, error) {
	for _, path := range discordSocketPaths() {
		conn, err := net.Dial("unix", path)
		if err == nil {
			return &ipcClient{conn: conn}, nil
		}
	}
	return nil, fmt.Errorf("discord: no IPC socket found")
}

func discordSocketPaths() []string {
	dirs := []string{
		os.Getenv("XDG_RUNTIME_DIR"),
		os.Getenv("TMPDIR"),
		os.TempDir(),
	}
	seen := make(map[string]bool)
	var paths []string
	for _, dir := range dirs {
		if dir == "" || seen[dir] {
			continue
		}
		seen[dir] = true
		for i := 0; i < 10; i++ {
			paths = append(paths, filepath.Join(dir, fmt.Sprintf("discord-ipc-%d", i)))
		}
	}
	return paths
}

func (c *ipcClient) writeFrame(opcode uint32, v any) error {
	data, err := json.Marshal(v)
	if err != nil {
		return fmt.Errorf("discord: marshal: %w", err)
	}
	buf := make([]byte, 8+len(data))
	binary.LittleEndian.PutUint32(buf[0:4], opcode)
	binary.LittleEndian.PutUint32(buf[4:8], uint32(len(data)))
	copy(buf[8:], data)
	if _, err := c.conn.Write(buf); err != nil {
		return fmt.Errorf("discord: write: %w", err)
	}
	return nil
}

func (c *ipcClient) readFrame() (uint32, json.RawMessage, error) {
	header := make([]byte, 8)
	if _, err := io.ReadFull(c.conn, header); err != nil {
		return 0, nil, fmt.Errorf("discord: read header: %w", err)
	}
	opcode := binary.LittleEndian.Uint32(header[0:4])
	size := binary.LittleEndian.Uint32(header[4:8])
	payload := make([]byte, size)
	if _, err := io.ReadFull(c.conn, payload); err != nil {
		return 0, nil, fmt.Errorf("discord: read payload: %w", err)
	}
	return opcode, payload, nil
}

func (c *ipcClient) handshake(clientID string) error {
	msg := map[string]any{"v": 1, "client_id": clientID}
	if err := c.writeFrame(opcodeHandshake, msg); err != nil {
		return err
	}
	_, _, err := c.readFrame()
	return err
}

type activity struct {
	Type       int         `json:"type"`
	Details    string      `json:"details,omitempty"`
	State      string      `json:"state,omitempty"`
	Timestamps *timestamps `json:"timestamps,omitempty"`
	Assets     *assets     `json:"assets,omitempty"`
}

type timestamps struct {
	Start int64 `json:"start,omitempty"`
	End   int64 `json:"end,omitempty"`
}

type assets struct {
	LargeImage string `json:"large_image,omitempty"`
	LargeText  string `json:"large_text,omitempty"`
}

func (c *ipcClient) setActivity(pid int, a *activity) error {
	type setArgs struct {
		PID      int       `json:"pid"`
		Activity *activity `json:"activity"`
	}
	type cmd struct {
		Cmd   string  `json:"cmd"`
		Args  setArgs `json:"args"`
		Nonce string  `json:"nonce"`
	}
	msg := cmd{
		Cmd:   "SET_ACTIVITY",
		Args:  setArgs{PID: pid, Activity: a},
		Nonce: nextNonce(),
	}
	if err := c.writeFrame(opcodeFrame, msg); err != nil {
		return err
	}
	_, _, err := c.readFrame()
	return err
}

func (c *ipcClient) close() {
	_ = c.writeFrame(opcodeClose, map[string]any{})
	c.conn.Close()
}
