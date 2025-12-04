package model

import (
	"time"

	connect "github.com/1Password/connect-sdk-go/onepassword"
	sdk "github.com/1password/onepassword-sdk-go"
)

// ItemSection represents a section within a 1Password item.
type ItemSection struct {
	ID    string
	Title string
}

// Item represents 1Password item.
type Item struct {
	ID        string
	VaultID   string
	Version   int
	Tags      []string
	Fields    []ItemField
	Sections  []ItemSection
	Files     []File
	CreatedAt time.Time
}

// FromConnectItem populates the Item from a Connect item.
func (i *Item) FromConnectItem(item *connect.Item) {
	i.ID = item.ID
	i.VaultID = item.Vault.ID
	i.Version = item.Version

	i.Tags = append(i.Tags, item.Tags...)

	// Build sections map from field references (Connect doesn't have a top-level sections array)
	sectionMap := make(map[string]ItemSection)
	for _, field := range item.Fields {
		if field.Section != nil {
			sectionID := field.Section.ID
			if _, exists := sectionMap[sectionID]; !exists {
				sectionMap[sectionID] = ItemSection{
					ID:    field.Section.ID,
					Title: field.Section.Label,
				}
			}
		}
	}
	// Convert map to slice
	i.Sections = make([]ItemSection, 0, len(sectionMap))
	for _, section := range sectionMap {
		i.Sections = append(i.Sections, section)
	}

	// Preserve fields with section and type metadata
	for _, field := range item.Fields {
		sectionID := ""
		if field.Section != nil {
			sectionID = field.Section.ID
		}
		i.Fields = append(i.Fields, ItemField{
			ID:        field.ID,
			Label:     field.Label,
			Value:     field.Value,
			SectionID: sectionID,
			FieldType: string(field.Type),
		})
	}

	for _, file := range item.Files {
		i.Files = append(i.Files, File{
			ID:   file.ID,
			Name: file.Name,
			Size: file.Size,
		})
	}

	i.CreatedAt = item.CreatedAt
}

// FromSDKItem populates the Item from an SDK item.
func (i *Item) FromSDKItem(item *sdk.Item) {
	i.ID = item.ID
	i.VaultID = item.VaultID
	i.Version = int(item.Version)

	i.Tags = make([]string, len(item.Tags))
	copy(i.Tags, item.Tags)

	// Preserve sections
	i.Sections = make([]ItemSection, len(item.Sections))
	for idx, section := range item.Sections {
		i.Sections[idx] = ItemSection{
			ID:    section.ID,
			Title: section.Title,
		}
	}

	// Preserve fields with section and type metadata
	for _, field := range item.Fields {
		sectionID := ""
		if field.SectionID != nil {
			sectionID = *field.SectionID
		}
		i.Fields = append(i.Fields, ItemField{
			ID:        field.ID,
			Label:     field.Title,
			Value:     field.Value,
			SectionID: sectionID,
			FieldType: string(field.FieldType),
		})
	}

	for _, file := range item.Files {
		i.Files = append(i.Files, File{
			ID:   file.Attributes.ID,
			Name: file.Attributes.Name,
			Size: int(file.Attributes.Size),
		})
	}

	// Items of 'Document' category keeps file information in the Document field.
	if item.Category == sdk.ItemCategoryDocument {
		i.Files = append(i.Files, File{
			ID:   item.Document.ID,
			Name: item.Document.Name,
			Size: int(item.Document.Size),
		})
	}

	i.CreatedAt = item.CreatedAt
}

// FromSDKItemOverview populates the Item from an SDK item overview.
func (i *Item) FromSDKItemOverview(item *sdk.ItemOverview) {
	i.ID = item.ID
	i.VaultID = item.VaultID

	i.Tags = make([]string, len(item.Tags))
	copy(i.Tags, item.Tags)

	i.CreatedAt = item.CreatedAt
}
