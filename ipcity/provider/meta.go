package provider

import (
	"fmt"
	"strconv"
	"strings"
)

// Meta defines the meta row information.
type Meta struct {
	country     string
	province    string
	city        string
	district    string
	isp         string
	backboneISP string
	countryCode int
	areaCode    int
}

// NewMeta returns a new meta.
func NewMeta() *Meta {
	return &Meta{}
}

// WithCountry returns the meta with the country.
func (r *Meta) WithCountry(c string) *Meta {
	if r != nil {
		r.country = c
	}
	return r
}

// WithProvince returns the meta with the province.
func (r *Meta) WithProvince(p string) *Meta {
	if r != nil {
		r.province = p
	}
	return r
}

// WithCity returns the meta with the city.
func (r *Meta) WithCity(c string) *Meta {
	if r != nil {
		r.city = c
	}
	return r
}

// WithDistrict returns the meta with the district.
func (r *Meta) WithDistrict(d string) *Meta {
	if r != nil {
		r.district = d
	}
	return r
}

// WithISP returns the meta with the ISP.
func (r *Meta) WithISP(i string) *Meta {
	if r != nil {
		r.isp = i
	}
	return r
}

// WithSeniorISP returns the meta with the senior ISP.
func (r *Meta) WithSeniorISP(s string) *Meta {
	return r.WithBackboneISP(s)
}

// WithBackboneISP returns the meta with the backbone ISP.
func (r *Meta) WithBackboneISP(s string) *Meta {
	if r != nil {
		r.backboneISP = s
	}
	return r
}

// WithCountryCode returns the meta with the country code.
func (r *Meta) WithCountryCode(c int) *Meta {
	if r != nil {
		r.countryCode = c
	}
	return r
}

// WithAreaCode returns the meta with the area code.
func (r *Meta) WithAreaCode(a int) *Meta {
	if r != nil {
		r.areaCode = a
	}
	return r
}

// Country returns the Country in the meta row information.
func (r *Meta) Country() string {
	if r != nil {
		return r.country
	}
	return ""
}

// Province returns the Province in the meta row information.
func (r *Meta) Province() string {
	if r != nil {
		return r.province
	}
	return ""
}

// City returns the City in the meta row information.
func (r *Meta) City() string {
	if r != nil {
		return r.city
	}
	return ""
}

// District returns the District in the meta row information.
func (r *Meta) District() string {
	if r != nil {
		return r.district
	}
	return ""
}

// ISP returns the ISP in the meta row information.
func (r *Meta) ISP() string {
	if r != nil {
		return r.isp
	}
	return ""
}

// SeniorISP returns the BackboneISP in the meta row information.
func (r *Meta) SeniorISP() string {
	return r.BackboneISP()
}

// BackboneISP returns the BackboneISP in the meta row information.
func (r *Meta) BackboneISP() string {
	if r != nil {
		return r.backboneISP
	}
	return ""
}

// CountryCode returns the CountryCode in the meta row information.
func (r *Meta) CountryCode() int {
	if r != nil {
		return r.countryCode
	}
	return 0
}

// AreaCode returns the AreaCode in the meta row information.
func (r *Meta) AreaCode() int {
	if r != nil {
		return r.areaCode
	}
	return 0
}

// IsEmpty returns true if all the fields are empty.
func (r *Meta) IsEmpty() bool {
	return r.Country() == "" &&
		r.Province() == "" &&
		r.City() == "" &&
		r.District() == "" &&
		r.ISP() == "" &&
		r.BackboneISP() == "" &&
		r.CountryCode() == 0 &&
		r.AreaCode() == 0
}

func (r *Meta) String() string {
	if r != nil {
		return fmt.Sprintf("country:%q province:%q city:%q district:%q "+
			"ISP:%q backboneISP:%q countryCode:%d areaCode:%d",
			r.Country(), r.Province(), r.City(), r.District(),
			r.ISP(), r.BackboneISP(), r.CountryCode(), r.AreaCode())
	}
	return ""
}

// UnmarshalString will fill the details into meta row.
func (r *Meta) UnmarshalString(line string) error {
	toInt := func(s string) (int, error) {
		if s == "" {
			return 0, nil
		}
		return strconv.Atoi(s)
	}

	var e error
	if r != nil {
		for i, item := range strings.Split(
			strings.TrimSuffix(line, "\n"), "\t") {
			switch i {
			case 0:
				r.country = item
			case 1:
				r.province = item
			case 2:
				r.city = item
			case 3:
				r.district = item
			case 4:
				r.isp = item
			case 5:
				r.backboneISP = item
			case 6:
				if r.countryCode, e = toInt(item); e != nil {
					break
				}
			case 7:
				if r.areaCode, e = toInt(item); e != nil {
					break
				}
			}
		}
	}
	return e
}

// Unmarshal a bytes array and fill the details into meta row.
func (r *Meta) Unmarshal(buffer []byte) error {
	return r.UnmarshalString(string(buffer[:]))
}

// MarshalString will serialize the data entity to a string.
func (r *Meta) MarshalString() (string, error) {
	var line string
	if r != nil {
		line = fmt.Sprintf("%s\t%s\t%s\t%s\t%s\t%s\t%d\t%d",
			r.Country(), r.Province(), r.City(), r.District(),
			r.ISP(), r.BackboneISP(), r.CountryCode(), r.AreaCode())
	}
	return line, nil
}

// Marshal will serialize the data entity to a string.
func (r *Meta) Marshal() ([]byte, error) {
	line, err := r.MarshalString()
	return []byte(line), err
}
