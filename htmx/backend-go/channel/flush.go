package channel

import "io"

func Flush(w io.Writer, ch chan []byte, delim ...byte) error {
	for b := range ch {
		_, err := w.Write(b)
		if err != nil {
			return err
		}
		if len(delim) != 0 {
			_, err := w.Write(delim)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
