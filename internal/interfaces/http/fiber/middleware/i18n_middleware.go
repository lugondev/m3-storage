package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

const translatorContextKey = "localizer"

// I18nMiddleware creates a Fiber middleware that detects the user's preferred language
// and makes a i18n.Translator available in the context.
func I18nMiddleware(bundle *i18n.Bundle) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get the Accept-Language header
		langHeader := c.Get("Accept-Language")

		// Match the preferred language
		// You can customize the matching behavior here
		matcher := language.NewMatcher(bundle.LanguageTags())
		tag, _ := language.MatchStrings(matcher, langHeader)

		// Create a translator for the matched language
		translator := i18n.NewLocalizer(bundle, tag.String())

		// Store the translator in the context
		c.Locals(translatorContextKey, translator)

		// Continue to the next middleware/handler
		return c.Next()
	}
}

// GetTranslator retrieves the i18n.Localizer from the Fiber context.
func GetTranslator(c *fiber.Ctx) *i18n.Localizer {
	if translator, ok := c.Locals(translatorContextKey).(*i18n.Localizer); ok {
		return translator
	}
	return nil
}

func TranslatorTranslate(c *fiber.Ctx, args ...string) string {
	if len(args) == 0 {
		return "" // Return empty string if no arguments are provided
	}
	key := args[0]
	defaultTranslation := args[0]
	if len(args) > 1 {
		defaultTranslation = args[1]
	}
	translator := GetTranslator(c)
	if translator == nil {
		return defaultTranslation // Fallback to default translation if translator is not available
	}
	translated, err := translator.LocalizeMessage(&i18n.Message{
		ID: key,
	})
	if err != nil {
		return defaultTranslation // Fallback to the key if translation fails
	}
	return translated
}
