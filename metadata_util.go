package phonenumber

import "github.com/nyaruka/phonenumbers/v2"

func mergeNumberFormats(dst, src *phonenumbers.NumberFormat) {
	if src.Pattern != nil {
		dst.Pattern = src.Pattern
	}
	if src.Format != nil {
		dst.Format = src.Format
	}
	dst.LeadingDigitsPattern = append([]string{}, src.LeadingDigitsPattern...)
	if src.NationalPrefixFormattingRule != nil {
		dst.NationalPrefixFormattingRule = src.NationalPrefixFormattingRule
	}
	if src.DomesticCarrierCodeFormattingRule != nil {
		dst.DomesticCarrierCodeFormattingRule = src.DomesticCarrierCodeFormattingRule
	}
	if src.NationalPrefixOptionalWhenFormatting != nil {
		dst.NationalPrefixOptionalWhenFormatting = src.NationalPrefixOptionalWhenFormatting
	}
}
