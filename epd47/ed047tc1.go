package epd47

// Internal config register layout mirrors ed047tc1.c.
type reg struct {
	epLatchEnable   bool
	powerDisable    bool
	posPowerEnable  bool
	negPowerEnable  bool
	epSTV           bool
	epScanDirection bool
	epMode          bool
	epOutputEnable  bool
}

func (d *Device) pushCfgBit(bit bool) {
	d.bus.cfgClk(false)
	d.bus.cfgData(bit)
	d.bus.cfgClk(true)
}

func (d *Device) pushCfg() {
	d.bus.cfgStr(false)
	// reverse order like C push_cfg
	d.pushCfgBit(d.cfg.epOutputEnable)
	d.pushCfgBit(d.cfg.epMode)
	d.pushCfgBit(d.cfg.epScanDirection)
	d.pushCfgBit(d.cfg.epSTV)
	d.pushCfgBit(d.cfg.negPowerEnable)
	d.pushCfgBit(d.cfg.posPowerEnable)
	d.pushCfgBit(d.cfg.powerDisable)
	d.pushCfgBit(d.cfg.epLatchEnable)
	d.bus.cfgStr(true)
}

// pulseCKV in microseconds
func (d *Device) pulseCKV(highUS, lowUS int) {
	if highUS > 0 {
		d.bus.ckv(true)
		d.bus.sleepUS(highUS)
	}
	d.bus.ckv(false)
	if lowUS > 0 {
		d.bus.sleepUS(lowUS)
	}
}

// Power sequences
func (d *Device) PowerOn() {
	d.cfg.epScanDirection = true
	d.cfg.powerDisable = false
	d.pushCfg()
	d.bus.sleepUS(100_000)

	d.cfg.negPowerEnable = true
	d.pushCfg()
	d.bus.sleepUS(500_000)

	d.cfg.posPowerEnable = true
	d.pushCfg()
	d.bus.sleepUS(100_000)

	d.cfg.epSTV = true
	d.pushCfg()

	d.bus.sth(true) // input enable
}

func (d *Device) PowerOff() {
	d.cfg.posPowerEnable = false
	d.pushCfg()
	d.bus.sleepUS(10_000)

	d.cfg.negPowerEnable = false
	d.pushCfg()
	d.bus.sleepUS(100_000)

	d.cfg.powerDisable = true
	d.pushCfg()

	d.cfg.epSTV = false
	d.pushCfg()
}

func (d *Device) PowerOffAll() {
	d.cfg = reg{} // all false
	d.pushCfg()
}

// Frame control
func (d *Device) StartFrame() {
	d.cfg.epMode = true
	d.pushCfg()

	d.pulseCKV(1, 1)
	d.cfg.epSTV = false
	d.pushCfg()
	d.bus.sleepUS(1_000) // coarse busy delay
	d.pulseCKV(10, 10)
	d.cfg.epSTV = true
	d.pushCfg()
	d.pulseCKV(0, 10)

	d.cfg.epOutputEnable = true
	d.pushCfg()
	d.pulseCKV(1, 1)
}

func (d *Device) EndFrame() {
	d.cfg.epOutputEnable = false
	d.pushCfg()
	d.cfg.epMode = false
	d.pushCfg()
	d.pulseCKV(1, 1)
	d.pulseCKV(1, 1)
}

// Row helpers
func (d *Device) latchRow() {
	d.cfg.epLatchEnable = true
	d.pushCfg()
	d.cfg.epLatchEnable = false
	d.pushCfg()
}

// SkipRow uses approx timing from the C driver (ticks -> us heuristic).
func (d *Device) SkipRow() {
	d.pulseCKV(45, 5)
}
