package sound

/*
Audio Interface
If you want to implement audio playback yourself, you should follow this interface.
More specifically, set a function that returns the implementation of this interface to Cubism.LoadSound.
*/
type Sound interface {
	// Play the sound
	Play() error
	// Stop the sound
	Close()
}
