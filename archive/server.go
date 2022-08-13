package archive

import (
	"encoding/binary"
	"io"
	"net/http"
	"time"
)

const sharedIndex string = "shared-index.dat"
const importCgi string = "import.cgi"

type Server struct {
	*Archive
	srv *http.Server
}

type writer struct {
	w   io.Writer
	err error
}

func (w *writer) write(data interface{}) {
	if w.err == nil {
		w.err = binary.Write(w.w, binary.BigEndian, data)
	}
}

func NewServer() *Server {
	s := &Server{}
	srv := &http.Server{
		Addr:           ":6667",
		Handler:        http.HandlerFunc(s.incomingHandler),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	s.srv = srv
	s.srv.ListenAndServe()
	return &Server{}
}

func (s Server) BuildSharedIndex() {
}

func (s Server) incomingHandler(w http.ResponseWriter, req *http.Request) {
	if req.URL.Path == "/"+sharedIndex {
		w.WriteHeader(http.StatusMethodNotAllowed)
		io.WriteString(w, "Pulling from this archive server is unimplemented!\n")
		return
	}
	if req.URL.Path == "/"+importCgi && req.Method == http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		io.WriteString(w, "Pushing from this archive server is unimplemented!\n")
		return
	}
	w.WriteHeader(http.StatusNotFound)
}

func (s *Server) Write(output io.Writer) error {
	w := writer{w: output}

	// Write ArchiveFlags (unimplemented)
	w.write(&s.ArchiveFlags)

	// Write the admin channel
	w.write(&s.AdminChannel)

	// Write the number of alternate URIs
	w.write(&s.NumAltURIs)

	// Populate AltURIs with other known archive servers
	for i := 0; i < int(s.NumAltURIs); i++ {
		// TODO: why do we write the length here and not below?
		var length uint16
		w.write(&length)
		//		uri := make([]byte, int(length))
		w.write(s.AltURIs[i])
	}

	// Count the number of channels
	w.write(&s.NumChannels)

	// Write the channel hashes
	for i := 0; i < int(s.NumChannels); i++ {
		w.write(&s.ChannelHashes[i])
	}

	w.write(&s.NumMessages)

	// Write messages and append urls
	for i := 0; i < int(s.NumMessages); i++ {
		w.write(&s.Messages[i])
	}
	if w.err != nil {
		return w.err
	}
	return nil
}
