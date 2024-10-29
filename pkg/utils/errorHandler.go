package utils

import "log"

func ErrorHandler(err error, message string) error {
    if err != nil {
        log.Printf(message + ": %v", err)
        return err
    }

    return nil
}
