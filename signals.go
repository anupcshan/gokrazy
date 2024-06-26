package gokrazy

import (
	"log"
	"os"
	"os/signal"

	"golang.org/x/sys/unix"
)

func listenForSignals(sighup func()) {
	{
		c := make(chan os.Signal, 1)
		signal.Notify(c, unix.SIGHUP)
		go func() {
			for range c {
				sighup()
			}
		}()
	}

	{
		c := make(chan os.Signal, 1)
		signal.Notify(c, unix.SIGTERM, unix.SIGUSR1, unix.SIGUSR2)
		go func() {
			for sig := range c {
				log.Printf("received signal %s, killing services", sig)

				killSupervisedServicesAndUmountPerm(defaultSignalDelay)

				switch sig {
				case unix.SIGTERM:
					reboot(true)

				case unix.SIGUSR1:
					log.Println("Halting")
					code := uint(unix.LINUX_REBOOT_CMD_HALT)
					if err := unix.Reboot(int(code)); err != nil {
						log.Printf("HALT: %v", err)
					}

				case unix.SIGUSR2:
					log.Println("Powering off")
					if err := unix.Reboot(unix.LINUX_REBOOT_CMD_POWER_OFF); err != nil {
						log.Printf("POWER_OFF: %v", err)
					}

				}
			}
		}()
	}
}
