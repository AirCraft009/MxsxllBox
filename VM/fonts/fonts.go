package fonts

type Character struct {
	Bitmap  [8]byte
	ByteVal byte
}

type Font struct {
	name       string
	Characters []*Character
}

func newFont(name string) *Font {
	return &Font{
		name: name,
	}
}

func newCharacter(bitmap [8]byte, byteVal byte) *Character {
	return &Character{
		Bitmap:  bitmap,
		ByteVal: byteVal,
	}
}
