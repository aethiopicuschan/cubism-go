package sound

/*
音声のインターフェース
音声再生を自分で実装したい場合はこのインターフェースに従えば良い
より具体的には、このインターフェースの実装を返す関数をCubism.LoadSoundに設定する
*/
type Sound interface {
	// 音声を再生する
	Play()
	// 音声を停止する
	Close()
}
