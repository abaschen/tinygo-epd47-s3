package epd47

// writeByte drives D0..D7 then strobes CKH.
func (d *Device) writeByte(b byte) {
	// Drive only configured pins; nil-safe because we guard by mask
	m := d.dataMask
	// Bit 0 .. 7
	if (m & 0x01) != 0 {
		d.bus.dataPins[0]((b & 0x01) != 0)
	}
	if (m & 0x02) != 0 {
		d.bus.dataPins[1]((b & 0x02) != 0)
	}
	if (m & 0x04) != 0 {
		d.bus.dataPins[2]((b & 0x04) != 0)
	}
	if (m & 0x08) != 0 {
		d.bus.dataPins[3]((b & 0x08) != 0)
	}
	if (m & 0x10) != 0 {
		d.bus.dataPins[4]((b & 0x10) != 0)
	}
	if (m & 0x20) != 0 {
		d.bus.dataPins[5]((b & 0x20) != 0)
	}
	if (m & 0x40) != 0 {
		d.bus.dataPins[6]((b & 0x40) != 0)
	}
	if (m & 0x80) != 0 {
		d.bus.dataPins[7]((b & 0x80) != 0)
	}

	d.bus.ckh(true)
	d.bus.sleepUS(0)
	d.bus.ckh(false)
}

// writeLineBytes emits raw bytes quickly.
func (d *Device) writeLineBytes(buf []byte) {
	for i := 0; i < len(buf); i++ {
		d.writeByte(buf[i])
	}
}

// 1bpp line writer: caller pre-fills line1b.
func (d *Device) outputRow1bpp(lineLenBytes int, pulseHighUS int) {
	d.latchRow()
	if pulseHighUS <= 0 {
		pulseHighUS = 10
	}
	d.pulseCKV(pulseHighUS, 50)
	d.writeLineBytes(d.line1b[:lineLenBytes])
	d.pulseCKV(1, 1)
}
