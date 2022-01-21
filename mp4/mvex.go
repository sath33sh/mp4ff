package mp4

import (
	"io"

	"github.com/edgeware/mp4ff/bits"
)

// MvexBox - MovieExtendsBox (mevx)
//
// Contained in : Movie Box (moov)
//
// Its presence signals a fragmented asset
type MvexBox struct {
	Mehd     *MehdBox
	Trex     *TrexBox
	Trexs    []*TrexBox
	Children []Box
}

// NewMvexBox - Generate a new empty mvex box
func NewMvexBox() *MvexBox {
	return &MvexBox{}
}

// AddChild - Add a child box
func (m *MvexBox) AddChild(box Box) {

	switch box.Type() {
	case "mehd":
		m.Mehd = box.(*MehdBox)
	case "trex":
		if m.Trex == nil {
			m.Trex = box.(*TrexBox)
		}
		m.Trexs = append(m.Trexs, box.(*TrexBox))
	}
	m.Children = append(m.Children, box)
}

// DecodeMvex - box-specific decode
func DecodeMvex(hdr boxHeader, startPos uint64, r io.Reader) (Box, error) {
	children, err := DecodeContainerChildren(hdr, startPos+8, startPos+hdr.size, r)
	if err != nil {
		return nil, err
	}
	m := NewMvexBox()
	for _, c := range children {
		m.AddChild(c)
	}
	return m, nil
}

// DecodeMvex - box-specific decode
func DecodeMvexSR(hdr boxHeader, startPos uint64, sr bits.SliceReader) (Box, error) {
	children, err := DecodeContainerChildrenSR(hdr, startPos+8, startPos+hdr.size, sr)
	if err != nil {
		return nil, err
	}
	m := NewMvexBox()
	for _, c := range children {
		m.AddChild(c)
	}
	return m, nil
}

// Type - return box type
func (m *MvexBox) Type() string {
	return "mvex"
}

// Size - return calculated size
func (m *MvexBox) Size() uint64 {
	return containerSize(m.Children)
}

// GetChildren - list of child boxes
func (m *MvexBox) GetChildren() []Box {
	return m.Children
}

// Encode - write mvex container to w
func (m *MvexBox) Encode(w io.Writer) error {
	return EncodeContainer(m, w)
}

// Encode - write mvex container to sw
func (m *MvexBox) EncodeSW(sw bits.SliceWriter) error {
	return EncodeContainerSW(m, sw)
}

// Info - write box-specific information
func (m *MvexBox) Info(w io.Writer, specificBoxLevels, indent, indentStep string) error {
	return ContainerInfo(m, w, specificBoxLevels, indent, indentStep)
}

// GetTrex - get trex box for trackID
func (m *MvexBox) GetTrex(trackID uint32) (trex *TrexBox, ok bool) {
	for _, trex := range m.Trexs {
		if trex.TrackID == trackID {
			return trex, false
		}
	}
	return nil, true
}
