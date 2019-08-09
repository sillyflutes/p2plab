// Copyright 2019 Netflix, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/Netflix/p2plab/peer"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func init() {
	// UNIX Time is faster and smaller than most timestamps. If you set
	// zerolog.TimeFieldFormat to an empty string, logs will write with UNIX
	// time.
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
}

func main() {
	err := run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "ociget: %s\n", err)
		os.Exit(1)
	}
}

func run() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	p, err := peer.NewPeer(ctx, "./tmp/ociget")
	if err != nil {
		return err
	}

	f, err := os.Open("input")
	if err != nil {
		return err
	}

	log.Info().Msg("Adding 'input'")
	n, err := p.Add(ctx, f)
	if err != nil {
		return err
	}

	log.Info().Msgf("Getting 'input': %q", n.Cid())
	r, err := p.Get(ctx, n.Cid())
	if err != nil {
		return err
	}
	defer r.Close()

	content, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}

	log.Info().Msg("Writing 'output'")
	err = ioutil.WriteFile("output", content, 0644)
	if err != nil {
		return err
	}

	return nil
}