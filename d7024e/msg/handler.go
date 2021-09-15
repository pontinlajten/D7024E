package msg

import (
	d "main/d7024e"
)

func msgHandler(channel chan []byte, me *d.Contact, network *d.Network) {
	r_data := <-channel

}
