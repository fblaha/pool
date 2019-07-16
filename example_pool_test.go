package pool_test

import (
	"crypto/sha256"
	"crypto/sha512"
	"fmt"
	"hash"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/fblaha/pool"
)

// hashFactory factory function for hashes
type hashFactory func() hash.Hash

// hashOutput contains a hashing output or an error
type hashOutput struct {
	file   string
	length int
	value  []byte
	err    error
}

// hashWorkFactory creates a hash worker function for given file and given hashing algorithm
func hashWorkFactory(file string, hashFactory hashFactory, out chan<- hashOutput) pool.WorkerFunc {
	return func() {
		h := hashFactory()
		output := hashOutput{file: file, length: h.Size()}
		f, err := os.Open(file)
		if err != nil {
			// propagates an error to the output and returns
			output.err = err
			out <- output
			return
		}
		defer f.Close()

		// creates an instance of hash
		if _, err := io.Copy(h, f); err != nil {
			// propagates an error to the output and returns
			output.err = err
			out <- output
			return
		}
		// writes computed a hash to the output channel
		output.value = h.Sum(nil)
		out <- output

	}
}

// Example demonstrates usage of the pool executor with the graceful shutdown and the result/error propagation,
// The example iterates all files in the testdata directory.
// For each file computes sha256, sha384 and sha512 hash.
func Example() {
	// the output channel
	output := make(chan hashOutput)
	// slice with hash factories for each type of hash
	hashFactories := []hashFactory{sha256.New, sha512.New384, sha512.New}

	go func() {
		// constructs an executor with a single goroutine (only to have deterministic output in console)
		// a real world usage will probably use a higher number
		executor := pool.NewExecutor(1)
		testDir := filepath.Join("testdata")

		// read files from the testdata dir
		files, _ := ioutil.ReadDir(testDir)
		for _, file := range files {
			for _, factory := range hashFactories {
				// submit a hash work for an execution
				// for each file and each hash
				file := filepath.Join(testDir, file.Name())
				executor.SubmitFunc(hashWorkFactory(file, factory, output))
			}
		}
		// wait for the completion and shutdowns executor
		executor.ShutdownGracefully()
		// closes the output channel
		close(output)

	}()

	// iterates the output channel and prints the results
	for result := range output {
		if result.err != nil {
			fmt.Printf("hashing of %s failed: %+v\n", result.file, result.err)
			continue
		}
		fmt.Printf("%s: %x (hash length: %d bytes)\n", filepath.Base(result.file), result.value, result.length)
	}
	// Output:
	//1.txt: bf41cf94047f1a3443ca654a235bc8f830f7997da9b6f3b2b041a866bc6e3b6f (hash length: 32 bytes)
	//1.txt: 2aa9a9c2e8e0a4473812799fe31214cf6cee4e331e5f493d849a85ac0ebbb86fe3ef336b9b257d1fc635809071dda1fa (hash length: 48 bytes)
	//1.txt: a0109048ea5c5c8db36ddd573ecf3a3830e53773af979ce4ed2287f2970c8f29a2892244caf890217782e7fe39c5b94036ab9ecde3ddbe056ca36abf4df7ac66 (hash length: 64 bytes)
	//2.txt: 54811cbc6c86311729b0a33e26c89087881b36b9ca3217d15cb5196e35f9a7e3 (hash length: 32 bytes)
	//2.txt: dfe72d43d92f735ee211d162e204354ff71c2eec4236e04f5b0ccde348e964f68f7334deab350a9fab7d592dd2dfa0d8 (hash length: 48 bytes)
	//2.txt: 3641e308bbf0cb0c59420e84e30a0091e06f32e0d48f5aac60223cce7d0e3326c3d29755f2c26b44374c815cffdaf6dc579da441bd801a4c3cf9f6a14fcc3537 (hash length: 64 bytes)
	//3.txt: c03d93310d14ce82b3d8ce9ddc4e9a1ddd791271a88ae8151db6abee6ff98f6d (hash length: 32 bytes)
	//3.txt: 4ee28c5f7be28c00712410e18658c3681c42966c5f9dc174896ad3871b74dd543deaeffed27892aa2b0b297c60f39d9c (hash length: 48 bytes)
	//3.txt: 867216026210aebb29f434f6117ddbbf564551e2cea93632e2f2ef2651231ec2bd51c813c1777c1c1c351507f31827d64be03d25a62f08f00dd8f86cb19f6254 (hash length: 64 bytes)
}
