package model

import "fmt"

type uploadCriteria struct {
   Always string
   Eval   string
}

var UploadCriteria = &uploadCriteria{
  Always: "always",
  Eval: "eval",
}

func ValidateUploadCriteria(uc string) error {
  switch uc {
  case UploadCriteria.Always, UploadCriteria.Eval:
    return nil
  default:
    return fmt.Errorf("model upload_criteria value \"%s\" not supported", uc)
  }
}