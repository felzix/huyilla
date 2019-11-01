package engine

import (
	"bytes"
	"fmt"
	"github.com/felzix/huyilla/constants"
	C "github.com/felzix/huyilla/constants"
	"github.com/felzix/huyilla/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)


func TestWeb(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Web Suite")
}


var _ = Describe("Web", func() {
	Describe("http api", func() {
		NAME := "arana"
		PASS := "murakami"

		auth, _ := (&types.Auth{Name: NAME, Password: []byte(PASS)}).Marshal()

		var engine *Engine
		var token string

		BeforeSuite(func() {
			h, err := NewEngine(C.SEED, NewLakeWorldGenerator(3), NewMemoryDatabase())
			engine = h
			Expect(err).To(BeNil())
		})

		AfterSuite(func() {
			err := engine.World.WipeDatabase()
			Expect(err).To(BeNil())
		})

		Describe("ping", func() {
			It("ping returns pong", func() {
				res := requesty("GET", "/ping", nil, engine, nil)

				Expect(res.Code).To(Equal(http.StatusOK))
				body, err := ioutil.ReadAll(res.Body)
				Expect(err).To(BeNil())
				Expect(body).To(Equal([]byte("pong")))
			})
		})

		Describe("signup flow", func() {
			It("Doesn't yet exist", func() {
				res := requesty("GET", "/auth/exists/"+NAME, nil, engine, nil)
				body, err := ioutil.ReadAll(res.Body)
				Expect(err).To(BeNil())
				Expect(body).To(Equal([]byte("false")))
				Expect(res.Code).To(Equal(http.StatusOK))
			})

			It("Signs up", func() {
				res := requesty("POST", "/auth/signup", bytes.NewReader(auth), engine, map[string]string{
					"contentType": "application/protobuf",
				})

				body, err := ioutil.ReadAll(res.Body)
				Expect(err).To(BeNil())
				Expect(body).To(Equal([]byte("Signup successful!")))
				Expect(res.Code).To(Equal(http.StatusOK))

				player, err := engine.World.Player(NAME)
				Expect(err).To(BeNil())
				Expect(player).ToNot(BeNil())
				Expect(len(player.Token)).To(Equal(0))
				Expect(player.Name).To(Equal(NAME))
			})

			It("Now exists", func() {
				res := requesty("GET", "/auth/exists/"+NAME, nil, engine, nil)
				body, err := ioutil.ReadAll(res.Body)
				Expect(err).To(BeNil())
				Expect(body).To(Equal([]byte("true")))
				Expect(res.Code).To(Equal(http.StatusOK))
			})

			It("Logs in", func() {
				res := requesty("POST", "/auth/login", bytes.NewReader(auth), engine, map[string]string{
					"contentType": "application/protobuf",
				})

				body, err := ioutil.ReadAll(res.Body)
				Expect(err).To(BeNil())
				Expect(len(body) > 100).To(BeTrue(), fmt.Sprintf(`Body was too short: "%s"`, body))
				Expect(body[0]).To(Equal(byte('e')))
				Expect(res.Code).To(Equal(http.StatusOK))

				token = string(body)

				player, err := engine.World.Player(NAME)
				Expect(err).To(BeNil())
				Expect(player).ToNot(BeNil())
				Expect(len(player.Token) > 0).To(BeTrue(), "Player token is not set")
				Expect(player.Name).To(Equal(NAME))
			})

			It("Logs out", func() {
				res := requesty("POST", "/auth/logout", nil, engine, map[string]string{
					"contentType": "application/protobuf",
					"token":       token,
				})

				Expect(res.Code).To(Equal(http.StatusOK))
				body, err := ioutil.ReadAll(res.Body)
				Expect(err).To(BeNil())
				Expect(body).To(Equal([]byte("Logout successful!")))

				player, err := engine.World.Player(NAME)
				Expect(err).To(BeNil())
				Expect(player).ToNot(BeNil())
				Expect(len(player.Token)).To(Equal(0))
			})
		})

		Describe("get world age", func() {
			It("gets world age", func() {
				res := requesty("GET", "/world/age", nil, engine, nil)

				Expect(res.Code).To(Equal(http.StatusOK))
				body, err := ioutil.ReadAll(res.Body)
				Expect(err).To(BeNil())
				Expect(body).To(Equal([]byte("1")))
			})
		})

		Describe("get player", func() {
			It("Gets player info", func() {
				res := requesty("GET", "/world/player/"+NAME, nil, engine, map[string]string{
					"contentType": "application/protobuf",
					"token":       token,
				})

				Expect(res.Code).To(Equal(http.StatusOK))
				body, err := ioutil.ReadAll(res.Body)
				Expect(err).To(BeNil())
				var entity types.Entity
				err = entity.Unmarshal(body)
				Expect(err).To(BeNil())
				Expect(entity.PlayerName).To(Equal(NAME))
				Expect(entity.Id).ToNot(Equal(0))
			})
		})

		Describe("get chunk", func() {
			It("in range", func() {
				res := requesty("GET", fmt.Sprintf("/world/chunk/%d/%d/%d", 0, 0, constants.ACTIVE_CHUNK_RADIUS), nil, engine, map[string]string{
					"contentType": "application/protobuf",
					"token":       token,
				})

				Expect(res.Code).To(Equal(http.StatusOK))
				body, err := ioutil.ReadAll(res.Body)
				Expect(err).To(BeNil())

				var chunks types.Chunks
				err = chunks.Unmarshal(body)
				Expect(err).To(BeNil())
				Expect(len(chunks.Chunks)).To(Equal(1))
				Expect(chunks.Chunks[0]).ToNot(BeNil())
				Expect(len(chunks.Chunks[0].Voxels)).To(Equal(constants.CHUNK_LENGTH))
			})

			It("in range, radius=1", func() {
				res := requesty("GET", fmt.Sprintf("/world/chunk/%d/%d/%d?radius=1", 0, 0, constants.ACTIVE_CHUNK_RADIUS), nil, engine, map[string]string{
					"contentType": "application/protobuf",
					"token":       token,
				})

				Expect(res.Code).To(Equal(http.StatusOK))
				body, err := ioutil.ReadAll(res.Body)
				Expect(err).To(BeNil())

				var chunks types.Chunks
				err = chunks.Unmarshal(body)
				Expect(err).To(BeNil())
				Expect(len(chunks.Chunks)).To(Equal(27))
				Expect(chunks.Chunks[0]).ToNot(BeNil())
				Expect(chunks.Chunks[26]).ToNot(BeNil())
				Expect(len(chunks.Chunks[0].Voxels)).To(Equal(constants.CHUNK_LENGTH))
			})

			It("barely in range, radius is max", func() {
				res := requesty("GET", fmt.Sprintf("/world/chunk/%d/%d/%d?radius=%d", 0, 0, 0, constants.ACTIVE_CHUNK_RADIUS), nil, engine, map[string]string{
					"contentType": "application/protobuf",
					"token":       token,
				})

				Expect(res.Code).To(Equal(http.StatusOK))
				body, err := ioutil.ReadAll(res.Body)
				Expect(err).To(BeNil())

				var chunks types.Chunks
				err = chunks.Unmarshal(body)
				Expect(err).To(BeNil())
				Expect(len(chunks.Chunks)).To(Equal(7 * 7 * 7))
				Expect(chunks.Chunks[0]).ToNot(BeNil())
				Expect(chunks.Chunks[7*7*7-1]).ToNot(BeNil())
				Expect(len(chunks.Chunks[0].Voxels)).To(Equal(constants.CHUNK_LENGTH))
			})

			// The idea is that the initial readtime is so large that two of them don't fit into the per-test duration set on the commandline (15s)
			It("database caching works", func() {
				// fill cache
				requesty("GET", fmt.Sprintf("/world/chunk/%d/%d/%d?radius=3", 0, 0, constants.ACTIVE_CHUNK_RADIUS), nil, engine, map[string]string{
					"contentType": "application/protobuf",
					"token":       token,
				})

				// read from cache
				res := requesty("GET", fmt.Sprintf("/world/chunk/%d/%d/%d?radius=3", 0, 0, constants.ACTIVE_CHUNK_RADIUS), nil, engine, map[string]string{
					"contentType": "application/protobuf",
					"token":       token,
				})

				Expect(res.Code).To(Equal(http.StatusOK))
				body, err := ioutil.ReadAll(res.Body)
				Expect(err).To(BeNil())

				var chunks types.Chunks
				err = chunks.Unmarshal(body)
				Expect(err).To(BeNil())
				Expect(len(chunks.Chunks)).To(Equal(7 * 7 * 7))
				Expect(chunks.Chunks[0]).ToNot(BeNil())
				Expect(chunks.Chunks[7*7*7-1]).ToNot(BeNil())
				Expect(len(chunks.Chunks[0].Voxels)).To(Equal(constants.CHUNK_LENGTH))
			})

			It("out of range", func() {
				res := requesty("GET", fmt.Sprintf("/world/chunk/%d/%d/%d", 0, constants.ACTIVE_CHUNK_RADIUS+1, 0), nil, engine, map[string]string{
					"contentType": "application/protobuf",
					"token":       token,
				})

				Expect(res.Code).To(Equal(http.StatusForbidden))
				body, err := ioutil.ReadAll(res.Body)
				Expect(err).To(BeNil())
				Expect(string(body)).To(Equal("can only load nearby chunks\n"))
			})

			It("out of range: too big", func() {
				res := requesty("GET", fmt.Sprintf("/world/chunk/%d/%d/%d?radius=%d", 0, 0, 0, constants.ACTIVE_CHUNK_RADIUS+1), nil, engine, map[string]string{
					"contentType": "application/protobuf",
					"token":       token,
				})

				Expect(res.Code).To(Equal(http.StatusForbidden))
				body, err := ioutil.ReadAll(res.Body)
				Expect(err).To(BeNil())
				Expect(string(body)).To(Equal("can only load nearby chunks\n"))
			})

			It("out of range: max size and shifted", func() {
				res := requesty("GET", fmt.Sprintf("/world/chunk/%d/%d/%d?radius=%d", 1, 0, 0, constants.ACTIVE_CHUNK_RADIUS), nil, engine, map[string]string{
					"contentType": "application/protobuf",
					"token":       token,
				})

				Expect(res.Code).To(Equal(http.StatusForbidden))
				body, err := ioutil.ReadAll(res.Body)
				Expect(err).To(BeNil())
				Expect(string(body)).To(Equal("can only load nearby chunks\n"))
			})
		})
	})
})

func requesty(method, url string, body io.Reader, engine *Engine, headers map[string]string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, url, body)

	for key, value := range headers {
		req.Header.Add(key, value)
	}
	res := httptest.NewRecorder()
	Router(engine).ServeHTTP(res, req)

	return res
}
