package validator

import (
	"log"
	"strings"

	ut "github.com/go-playground/universal-translator"

	validatorGo "github.com/go-playground/validator/v10"
)

func registerCustomIndonesianTranslator(v *validatorGo.Validate, trans ut.Translator) {

	if err := v.RegisterTranslation("rfe", trans, func(ut ut.Translator) error {
		if err := ut.Add("rfe", "{0} wajib diisi jika {1} = {2}", false); err != nil {
			return err
		}
		return nil
	}, func(ut ut.Translator, fe validatorGo.FieldError) string {
		param := strings.Split(fe.Param(), `:`)
		paramField := param[0]
		paramValue := param[1]
		t, err := ut.T(fe.Tag(), fe.Field(), paramField, paramValue)
		if err != nil {
			log.Printf("warning: error translating FieldError: %#v", fe)
			return fe.(error).Error()
		}
		return t
	}); err != nil {
		panic(err)
	}

	if err := v.RegisterTranslation("enum", trans, func(ut ut.Translator) error {
		if err := ut.Add("enum", "kriteria yang diterima dari {0} adalah {1}", false); err != nil {
			return err
		}
		return nil
	}, func(ut ut.Translator, fe validatorGo.FieldError) string {
		// first, clean/remove the comma
		cleaned := strings.Replace(fe.Param(), "-", " ", -1)

		// convert 'cleaned' comma separated string to slice
		strSlice := strings.Fields(cleaned)

		t, err := ut.T(fe.Tag(), fe.Field(), strings.Join(strSlice, ","))
		if err != nil {
			log.Printf("warning: error translating FieldError: %#v", fe)
			return fe.(error).Error()
		}
		return t
	}); err != nil {
		panic(err)
	}

	if err := v.RegisterTranslation("date_only", trans, func(ut ut.Translator) error {
		if err := ut.Add("date_only", "{0} harus format yang benar yyyy-mm-dd.", false); err != nil {
			return err
		}
		return nil
	}, func(ut ut.Translator, fe validatorGo.FieldError) string {
		t, err := ut.T(fe.Tag(), fe.Field(), fe.Param())
		if err != nil {
			log.Printf("warning: error translating FieldError: %#v", fe)
			return fe.(error).Error()
		}
		return t
	}); err != nil {
		panic(err)
	}

	if err := v.RegisterTranslation("date_range", trans, func(ut ut.Translator) error {
		if err := ut.Add("date_range", "{0} harus berupa tanggal yang valid.", false); err != nil {
			return err
		}
		return nil
	}, func(ut ut.Translator, fe validatorGo.FieldError) string {
		t, err := ut.T(fe.Tag(), fe.Field(), fe.Param())
		if err != nil {
			log.Printf("warning: error translating FieldError: %#v", fe)
			return fe.(error).Error()
		}
		return t
	}); err != nil {
		panic(err)
	}

	if err := v.RegisterTranslation("no_space", trans, func(ut ut.Translator) error {
		if err := ut.Add("no_space", "field tidak bisa hanya berisi spasi.", false); err != nil {
			return err
		}
		return nil
	}, func(ut ut.Translator, fe validatorGo.FieldError) string {
		t, err := ut.T(fe.Tag(), fe.Field(), fe.Param())
		if err != nil {
			log.Printf("warning: error translating FieldError: %#v", fe)
			return fe.(error).Error()
		}
		return t
	}); err != nil {
		panic(err)
	}

	if err := v.RegisterTranslation("check_fk", trans, func(ut ut.Translator) error {
		if err := ut.Add("check_fk", "{0} tidak ditemukan.", false); err != nil {
			return err
		}
		return nil
	}, func(ut ut.Translator, fe validatorGo.FieldError) string {
		t, err := ut.T(fe.Tag(), fe.Field(), fe.Param())
		if err != nil {
			log.Printf("warning: error translating FieldError: %#v", fe)
			return fe.(error).Error()
		}
		return t
	}); err != nil {
		panic(err)
	}

	if err := v.RegisterTranslation("unique_store", trans, func(ut ut.Translator) error {
		if err := ut.Add("unique_store", "{0} sudah terdaftar.", false); err != nil {
			return err
		}
		return nil
	}, func(ut ut.Translator, fe validatorGo.FieldError) string {
		t, err := ut.T(fe.Tag(), fe.Field(), fe.Param())
		if err != nil {
			log.Printf("warning: error translating FieldError: %#v", fe)
			return fe.(error).Error()
		}
		return t
	}); err != nil {
		panic(err)
	}

	if err := v.RegisterTranslation("unique_update", trans, func(ut ut.Translator) error {
		if err := ut.Add("unique_update", "{0} sudah terdaftar.", false); err != nil {
			return err
		}
		return nil
	}, func(ut ut.Translator, fe validatorGo.FieldError) string {
		t, err := ut.T(fe.Tag(), fe.Field(), fe.Param())
		if err != nil {
			log.Printf("warning: error translating FieldError: %#v", fe)
			return fe.(error).Error()
		}
		return t
	}); err != nil {
		panic(err)
	}

}
