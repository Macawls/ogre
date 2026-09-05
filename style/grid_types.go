package style

type TrackSizingKind int

const (
	TrackFixed TrackSizingKind = iota
	TrackAuto
	TrackFr
)

type TrackSize struct {
	Kind   TrackSizingKind
	Length Value
	Fr     float64
}

type TrackList struct {
	Tracks []TrackSize
}

type GridAutoFlow int

const (
	FlowRow GridAutoFlow = iota
	FlowColumn
	FlowRowDense
	FlowColumnDense
)

type GridPlacement struct {
	Start int
	End   int
	Span  int
}
