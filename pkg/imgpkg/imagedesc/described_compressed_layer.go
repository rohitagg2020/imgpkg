// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

// Package imagedesc OCI Image descriptors
package imagedesc

import (
	"fmt"
	"io"

	regv1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/types"
	"github.com/vmware-tanzu/carvel-imgpkg/pkg/imgpkg/imageutils/gzip"
	"github.com/vmware-tanzu/carvel-imgpkg/pkg/imgpkg/imageutils/verify"
)

// DescribedCompressedLayer Represents a Layer that is part of a Tar file generated by imgpkg
type DescribedCompressedLayer struct {
	desc     ImageLayerDescriptor
	contents LayerContents
}

var _ regv1.Layer = DescribedCompressedLayer{}

// NewDescribedCompressedLayer Builds DescribedCompressedLayer struct
func NewDescribedCompressedLayer(desc ImageLayerDescriptor, contents LayerContents) DescribedCompressedLayer {
	return DescribedCompressedLayer{desc, contents}
}

// Digest returns the Digest of the layer
func (l DescribedCompressedLayer) Digest() (regv1.Hash, error) { return regv1.NewHash(l.desc.Digest) }

// DiffID returns the DiffID of the layer
func (l DescribedCompressedLayer) DiffID() (regv1.Hash, error) { return regv1.NewHash(l.desc.DiffID) }

// Compressed returns a reader for the Layer anv validates the Digest of the layer matches
func (l DescribedCompressedLayer) Compressed() (io.ReadCloser, error) {
	rc, err := l.contents.Open()
	if err != nil {
		return nil, err
	}

	h, err := l.Digest()
	if err != nil {
		return nil, fmt.Errorf("Computing digest: %v", err)
	}

	rc, err = verify.ReadCloser(rc, verify.SizeUnknown, h)
	if err != nil {
		return nil, fmt.Errorf("Creating verified reader: %v", err)
	}

	return rc, nil
}

// Uncompressed returns a reader for the Layer uncompressed
func (l DescribedCompressedLayer) Uncompressed() (io.ReadCloser, error) {
	rc, err := l.Compressed()
	if err != nil {
		return nil, err
	}

	return gzip.UnzipReadCloser(rc)
}

// Size returns the size of the Layer
func (l DescribedCompressedLayer) Size() (int64, error) { return l.desc.Size, nil }

// MediaType returns the Media Type of the layer
func (l DescribedCompressedLayer) MediaType() (types.MediaType, error) {
	return types.MediaType(l.desc.MediaType), nil
}
