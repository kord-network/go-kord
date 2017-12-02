// This file is part of the go-meta library.
//
// Copyright (C) 2017 JAAK MUSIC LTD
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.
//
// If you have any questions please contact yo@jaak.io

package ern

type NewReleaseMessage struct {
	MessageHeader   *MessageHeader  `xml:"MessageHeader,omitempty,omitempty"`
	UpdateIndicator string          `xml:"UpdateIndicator,omitempty" json:"UpdateIndicator,omitempty"`
	IsBackfill      bool            `xml:"IsBackfill,omitempty" json:"IsBackfill,omitempty"`
	WorkList        *WorkList       `xml:"WorkList,omitempty" json:"WorkList,omitempty"`
	CueSheetList    *CueSheetList   `xml:"CueSheetList,omitempty" json:"CueSheetList,omitempty"`
	ResourceList    *ResourceList   `xml:"ResourceList,omitempty" json:"ResourceList,omitempty"`
	CollectionList  *CollectionList `xml:"CollectionList,omitempty" json:"CollectionList,omitempty"`
	ReleaseList     *ReleaseList    `xml:"ReleaseList,omitempty" json:"ReleaseList,omitempty"`
	DealList        *DealList       `xml:"DealList,omitempty" json:"DealList,omitempty"`
}

type MessageHeader struct {
	MessageThreadId        string             `xml:"MessageThreadId,omitempty" json:"MessageThreadId,omitempty"`
	MessageId              string             `xml:"MessageId,omitempty" json:"MessageId,omitempty"`
	MessageFileName        string             `xml:"MessageFileName,omitempty" json:"MessageFileName,omitempty"`
	MessageSender          *MessagingParty    `xml:"MessageSender,omitempty" json:"MessageSender,omitempty"`
	SentOnBehalfOf         *MessagingParty    `xml:"SentOnBehalfOf,omitempty" json:"SentOnBehalfOf,omitempty"`
	MessageRecipient       *MessagingParty    `xml:"MessageRecipient,omitempty" json:"MessageRecipient,omitempty"`
	MessageCreatedDateTime string             `xml:"MessageCreatedDateTime,omitempty" json:"MessageCreatedDateTime,omitempty"`
	MessageAuditTrail      *MessageAuditTrail `xml:"MessageAuditTrail,omitempty" json:"MessageAuditTrail,omitempty"`
	Comment                *Comment           `xml:"Comment,omitempty" json:"Comment,omitempty"`
	MessageControlType     string             `xml:"MessageControlType,omitempty" json:"MessageControlType,omitempty"`
	LanguageAndScriptCode  string             `xml:"LanguageAndScriptCode,attr,omitempty" json:"LanguageAndScriptCode,omitempty"`
}

type MessagingParty struct {
	PartyId     *PartyId   `xml:"PartyId,omitempty" json:"PartyId,omitempty"`
	PartyName   *PartyName `xml:"PartyName,omitempty" json:"PartyName,omitempty"`
	TradingName *Name      `xml:"TradingName,omitempty" json:"TradingName,omitempty"`
}

type PartyId struct {
	Value     string `xml:",chardata"`
	Namespace string `xml:"Namespace,attr,omitempty" json:"Namespace,omitempty"`
	IsDPID    bool   `xml:"IsDPID,attr,omitempty" json:"IsDPID,omitempty"`
	IsISNI    bool   `xml:"IsISNI,attr,omitempty" json:"IsISNI,omitempty"`
}

type PartyName struct {
	FullName              *Name  `xml:"FullName,omitempty" json:"FullName,omitempty"`
	LanguageAndScriptCode string `xml:"LanguageAndScriptCode,attr,omitempty" json:"LanguageAndScriptCode,omitempty"`
}

type Name struct {
	Value                 string `xml:",chardata"`
	LanguageAndScriptCode string `xml:"LanguageAndScriptCode,attr,omitempty" json:"LanguageAndScriptCode,omitempty"`
}

type MessageAuditTrail struct {
	MessageAuditTrailEvent *MessageAuditTrailEvent `xml:"MessageAuditTrailEvent,omitempty" json:"MessageAuditTrailEvent,omitempty"`
	LanguageAndScriptCode  string                  `xml:"LanguageAndScriptCode,attr,omitempty" json:"LanguageAndScriptCode,omitempty"`
}

type MessageAuditTrailEvent struct {
	MessagingPartyDescriptor *MessagingParty `xml:"MessagingPartyDescriptor,omitempty" json:"MessagingPartyDescriptor,omitempty"`
	DateTime                 string          `xml:"DateTime,omitempty" json:"DateTime,omitempty"`
}

type DealList struct {
	ReleaseDeal           []*ReleaseDeal `xml:"ReleaseDeal,omitempty" json:"ReleaseDeal,omitempty"`
	LanguageAndScriptCode string         `xml:"LanguageAndScriptCode,attr,omitempty" json:"LanguageAndScriptCode,omitempty"`
}

type ReleaseDeal struct {
	DealReleaseReference  string  `xml:"DealReleaseReference,omitempty" json:"DealReleaseReference,omitempty"`
	Deal                  []*Deal `xml:"Deal,omitempty" json:"Deal,omitempty"`
	EffectiveDate         string  `xml:"EffectiveDate,omitempty" json:"EffectiveDate,omitempty"`
	LanguageAndScriptCode string  `xml:"LanguageAndScriptCode,attr,omitempty" json:"LanguageAndScriptCode,omitempty"`
}

type Deal struct {
	DealReference                             *DealReference                             `xml:"DealReference,omitempty" json:"DealReference,omitempty"`
	DealTerms                                 *DealTerms                                 `xml:"DealTerms,omitempty" json:"DealTerms,omitempty"`
	ResourceUsage                             *ResourceUsage                             `xml:"ResourceUsage,omitempty" json:"ResourceUsage,omitempty"`
	DealTechnicalResourceDetailsReferenceList *DealTechnicalResourceDetailsReferenceList `xml:"DealTechnicalResourceDetailsReferenceList,omitempty" json:"DealTechnicalResourceDetailsReferenceList,omitempty"`
	DistributionChannelPage                   *WebPage                                   `xml:"DistributionChannelPage,omitempty" json:"DistributionChannelPage,omitempty"`
	LanguageAndScriptCode                     string                                     `xml:"LanguageAndScriptCode,attr,omitempty" json:"LanguageAndScriptCode,omitempty"`
}

type DealTechnicalResourceDetailsReferenceList struct {
	DealTechnicalResourceDetailsReference []string `xml:"DealTechnicalResourceDetailsReference,omitempty" json:"DealTechnicalResourceDetailsReference,omitempty"`
}

type ResourceUsage struct {
	DealResourceReference string `xml:"DealResourceReference,omitempty" json:"DealResourceReference,omitempty"`
	Usage                 *Usage `xml:"Usage,omitempty" json:"Usage,omitempty"`
}

type DealTerms struct {
	IsPreOrderDeal                   bool                       `xml:"IsPreOrderDeal,omitempty" json:"IsPreOrderDeal,omitempty"`
	CommercialModelType              *CommercialModelType       `xml:"CommercialModelType,omitempty" json:"CommercialModelType,omitempty"`
	Usage                            *Usage                     `xml:"Usage,omitempty" json:"Usage,omitempty"`
	AllDealsCancelled                bool                       `xml:"AllDealsCancelled,omitempty" json:"AllDealsCancelled,omitempty"`
	TakeDown                         bool                       `xml:"TakeDown,omitempty" json:"TakeDown,omitempty"`
	TerritoryCode                    *TerritoryCode             `xml:"TerritoryCode,omitempty" json:"TerritoryCode,omitempty"`
	ExcludedTerritoryCode            *TerritoryCode             `xml:"ExcludedTerritoryCode,omitempty" json:"ExcludedTerritoryCode,omitempty"`
	DistributionChannel              *DSP                       `xml:"DistributionChannel,omitempty" json:"DistributionChannel,omitempty"`
	ExcludedDistributionChannel      *DSP                       `xml:"ExcludedDistributionChannel,omitempty" json:"ExcludedDistributionChannel,omitempty"`
	PriceInformation                 *PriceInformation          `xml:"PriceInformation,omitempty" json:"PriceInformation,omitempty"`
	IsPromotional                    bool                       `xml:"IsPromotional,omitempty" json:"IsPromotional,omitempty"`
	PromotionalCode                  *PromotionalCode           `xml:"PromotionalCode,omitempty" json:"PromotionalCode,omitempty"`
	ValidityPeriod                   *Period                    `xml:"ValidityPeriod,omitempty" json:"ValidityPeriod,omitempty"`
	ConsumerRentalPeriod             *ConsumerRentalPeriod      `xml:"ConsumerRentalPeriod,omitempty" json:"ConsumerRentalPeriod,omitempty"`
	PreOrderReleaseDate              *EventDate                 `xml:"PreOrderReleaseDate,omitempty" json:"PreOrderReleaseDate,omitempty"`
	ReleaseDisplayStartDate          string                     `xml:"ReleaseDisplayStartDate,omitempty" json:"ReleaseDisplayStartDate,omitempty"`
	TrackListingPreviewStartDate     string                     `xml:"TrackListingPreviewStartDate,omitempty" json:"TrackListingPreviewStartDate,omitempty"`
	CoverArtPreviewStartDate         string                     `xml:"CoverArtPreviewStartDate,omitempty" json:"CoverArtPreviewStartDate,omitempty"`
	ClipPreviewStartDate             string                     `xml:"ClipPreviewStartDate,omitempty" json:"ClipPreviewStartDate,omitempty"`
	ReleaseDisplayStartDateTime      string                     `xml:"ReleaseDisplayStartDateTime,omitempty" json:"ReleaseDisplayStartDateTime,omitempty"`
	TrackListingPreviewStartDateTime string                     `xml:"TrackListingPreviewStartDateTime,omitempty" json:"TrackListingPreviewStartDateTime,omitempty"`
	CoverArtPreviewStartDateTime     string                     `xml:"CoverArtPreviewStartDateTime,omitempty" json:"CoverArtPreviewStartDateTime,omitempty"`
	ClipPreviewStartDateTime         string                     `xml:"ClipPreviewStartDateTime,omitempty" json:"ClipPreviewStartDateTime,omitempty"`
	PreOrderPreviewDate              *EventDate                 `xml:"PreOrderPreviewDate,omitempty" json:"PreOrderPreviewDate,omitempty"`
	PreOrderPreviewDateTime          string                     `xml:"PreOrderPreviewDateTime,omitempty" json:"PreOrderPreviewDateTime,omitempty"`
	PreOrderIncentiveResourceList    *DealResourceReferenceList `xml:"PreOrderIncentiveResourceList,omitempty" json:"PreOrderIncentiveResourceList,omitempty"`
	InstantGratificationResourceList *DealResourceReferenceList `xml:"InstantGratificationResourceList,omitempty" json:"InstantGratificationResourceList,omitempty"`
	IsExclusive                      bool                       `xml:"IsExclusive,omitempty" json:"IsExclusive,omitempty"`
	RelatedReleaseOfferSet           *RelatedReleaseOfferSet    `xml:"RelatedReleaseOfferSet,omitempty" json:"RelatedReleaseOfferSet,omitempty"`
	PhysicalReturns                  *PhysicalReturns           `xml:"PhysicalReturns,omitempty" json:"PhysicalReturns,omitempty"`
	NumberOfProductsPerCarton        int                        `xml:"NumberOfProductsPerCarton,omitempty" json:"NumberOfProductsPerCarton,omitempty"`
	RightsClaimPolicy                *RightsClaimPolicy         `xml:"RightsClaimPolicy,omitempty" json:"RightsClaimPolicy,omitempty"`
	WebPolicy                        *WebPolicy                 `xml:"WebPolicy,omitempty" json:"WebPolicy,omitempty"`
	LanguageAndScriptCode            string                     `xml:"LanguageAndScriptCode,attr,omitempty" json:"LanguageAndScriptCode,omitempty"`
}

type WebPolicy struct {
	Condition               *Condition `xml:"Condition,omitempty" json:"Condition,omitempty"`
	AccessBlockingRequested bool       `xml:"AccessBlockingRequested,omitempty" json:"AccessBlockingRequested,omitempty"`
	AccessLimitation        string     `xml:"AccessLimitation,omitempty" json:"AccessLimitation,omitempty"`
	EmbeddingAllowed        bool       `xml:"EmbeddingAllowed,omitempty" json:"EmbeddingAllowed,omitempty"`
	UserRatingAllowed       bool       `xml:"UserRatingAllowed,omitempty" json:"UserRatingAllowed,omitempty"`
	UserCommentAllowed      bool       `xml:"UserCommentAllowed,omitempty" json:"UserCommentAllowed,omitempty"`
	UserResponsesAllowed    bool       `xml:"UserResponsesAllowed,omitempty" json:"UserResponsesAllowed,omitempty"`
	SyndicationAllowed      bool       `xml:"SyndicationAllowed,omitempty" json:"SyndicationAllowed,omitempty"`
}

type RightsClaimPolicy struct {
	Condition             *Condition `xml:"Condition,omitempty" json:"Condition,omitempty"`
	RightsClaimPolicyType string     `xml:"RightsClaimPolicyType,omitempty" json:"RightsClaimPolicyType,omitempty"`
}

type Condition struct {
	Value             string `xml:"Value,omitempty" json:"Value,omitempty"`
	Unit              string `xml:"Unit,omitempty" json:"Unit,omitempty"`
	ReferenceCreation string `xml:"ReferenceCreation,omitempty" json:"ReferenceCreation,omitempty"`
	RelationalRelator string `xml:"RelationalRelator,omitempty" json:"RelationalRelator,omitempty"`
}

type PhysicalReturns struct {
	PhysicalReturnsAllowed       bool   `xml:"PhysicalReturnsAllowed,omitempty" json:"PhysicalReturnsAllowed,omitempty"`
	LatestDateForPhysicalReturns string `xml:"LatestDateForPhysicalReturns,omitempty" json:"LatestDateForPhysicalReturns,omitempty"`
}

type RelatedReleaseOfferSet struct {
	ReleaseId             *ReleaseId   `xml:"ReleaseId,omitempty" json:"ReleaseId,omitempty"`
	ReleaseDescription    *Description `xml:"ReleaseDescription,omitempty" json:"ReleaseDescription,omitempty"`
	Deal                  *Deal        `xml:"Deal,omitempty" json:"Deal,omitempty"`
	LanguageAndScriptCode string       `xml:"LanguageAndScriptCode,attr,omitempty" json:"LanguageAndScriptCode,omitempty"`
}

type DealResourceReferenceList struct {
	DealResourceReference []string `xml:"DealResourceReference,omitempty" json:"DealResourceReference,omitempty"`
	Period                *Period  `xml:"Period,omitempty" json:"Period,omitempty"`
}

type ConsumerRentalPeriod struct {
	Value        string `xml:",chardata"`
	IsExtensible bool   `xml:"IsExtensible,attr,omitempty" json:"IsExtensible,omitempty"`
}

type PromotionalCode struct {
	Value     string `xml:",chardata"`
	Namespace string `xml:"Namespace,attr,omitempty" json:"Namespace,omitempty"`
}

type PriceInformation struct {
	Description                    *Description    `xml:"Description,omitempty" json:"Description,omitempty"`
	PriceRangeType                 *PriceRangeType `xml:"PriceRangeType,omitempty" json:"PriceRangeType,omitempty"`
	PriceType                      *PriceType      `xml:"PriceType,omitempty" json:"PriceType,omitempty"`
	WholesalePricePerUnit          *Price          `xml:"WholesalePricePerUnit,omitempty" json:"WholesalePricePerUnit,omitempty"`
	BulkOrderWholesalePricePerUnit *Price          `xml:"BulkOrderWholesalePricePerUnit,omitempty" json:"BulkOrderWholesalePricePerUnit,omitempty"`
	SuggestedRetailPrice           *Price          `xml:"SuggestedRetailPrice,omitempty" json:"SuggestedRetailPrice,omitempty"`
	PriceTypeAttr                  string          `xml:"PriceType,attr,omitempty" json:"PriceTypeAttr,omitempty"`
}

type Price struct {
	Value        string `xml:",chardata"`
	CurrencyCode string `xml:"CurrencyCode,attr,omitempty" json:"CurrencyCode,omitempty"`
}

type PriceType struct {
	Value     string `xml:",chardata"`
	Namespace string `xml:"Namespace,attr,omitempty" json:"Namespace,omitempty"`
}

type PriceRangeType struct {
	Value     string `xml:",chardata"`
	Namespace string `xml:"Namespace,attr,omitempty" json:"Namespace,omitempty"`
}

type DSP struct {
	PartyId               *PartyId       `xml:"PartyId,omitempty" json:"PartyId,omitempty"`
	PartyName             *PartyName     `xml:"PartyName,omitempty" json:"PartyName,omitempty"`
	TradingName           *Name          `xml:"TradingName,omitempty" json:"TradingName,omitempty"`
	URL                   string         `xml:"URL,omitempty" json:"URL,omitempty"`
	TerritoryCode         *TerritoryCode `xml:"TerritoryCode,omitempty" json:"TerritoryCode,omitempty"`
	LanguageAndScriptCode string         `xml:"LanguageAndScriptCode,attr,omitempty" json:"LanguageAndScriptCode,omitempty"`
}

type Usage struct {
	UseType                 []*UseType               `xml:"UseType,omitempty" json:"UseType,omitempty"`
	UserInterfaceType       *UserInterfaceType       `xml:"UserInterfaceType,omitempty" json:"UserInterfaceType,omitempty"`
	DistributionChannelType *DistributionChannelType `xml:"DistributionChannelType,omitempty" json:"DistributionChannelType,omitempty"`
	CarrierType             *CarrierType             `xml:"CarrierType,omitempty" json:"CarrierType,omitempty"`
	TechnicalInstantiation  *TechnicalInstantiation  `xml:"TechnicalInstantiation,omitempty" json:"TechnicalInstantiation,omitempty"`
	NumberOfUsages          int                      `xml:"NumberOfUsages,omitempty" json:"NumberOfUsages,omitempty"`
}

type TechnicalInstantiation struct {
	DrmEnforcementType  string   `xml:"DrmEnforcementType,omitempty" json:"DrmEnforcementType,omitempty"`
	VideoDefinitionType string   `xml:"VideoDefinitionType,omitempty" json:"VideoDefinitionType,omitempty"`
	CodingType          string   `xml:"CodingType,omitempty" json:"CodingType,omitempty"`
	BitRate             *BitRate `xml:"BitRate,omitempty" json:"BitRate,omitempty"`
}

type DealReference struct {
	Value                 string `xml:",chardata"`
	LanguageAndScriptCode string `xml:"LanguageAndScriptCode,attr,omitempty" json:"LanguageAndScriptCode,omitempty"`
}

type ReleaseList struct {
	Release               []*Release `xml:"Release,omitempty" json:"Release,omitempty"`
	LanguageAndScriptCode string     `xml:"LanguageAndScriptCode,attr,omitempty" json:"LanguageAndScriptCode,omitempty"`
}

type Release struct {
	ReleaseId                      *ReleaseId                      `xml:"ReleaseId,omitempty" json:"ReleaseId,omitempty"`
	ReleaseReference               string                          `xml:"ReleaseReference,omitempty" json:"ReleaseReference,omitempty"`
	ExternalResourceLink           *ExternalResourceLink           `xml:"ExternalResourceLink,omitempty" json:"ExternalResourceLink,omitempty"`
	SalesReportingProxyReleaseId   *SalesReportingProxyReleaseId   `xml:"SalesReportingProxyReleaseId,omitempty" json:"SalesReportingProxyReleaseId,omitempty"`
	ReferenceTitle                 *ReferenceTitle                 `xml:"ReferenceTitle,omitempty" json:"ReferenceTitle,omitempty"`
	ReleaseResourceReferenceList   *ReleaseResourceReferenceList   `xml:"ReleaseResourceReferenceList,omitempty" json:"ReleaseResourceReferenceList,omitempty"`
	ResourceOmissionReason         *ResourceOmissionReason         `xml:"ResourceOmissionReason,omitempty" json:"ResourceOmissionReason,omitempty"`
	ReleaseCollectionReferenceList *ReleaseCollectionReferenceList `xml:"ReleaseCollectionReferenceList,omitempty" json:"ReleaseCollectionReferenceList,omitempty"`
	ReleaseType                    *ReleaseType                    `xml:"ReleaseType,omitempty" json:"ReleaseType,omitempty"`
	ReleaseDetailsByTerritory      []*ReleaseDetailsByTerritory    `xml:"ReleaseDetailsByTerritory,omitempty" json:"ReleaseDetailsByTerritory,omitempty"`
	LanguageOfPerformance          string                          `xml:"LanguageOfPerformance,omitempty" json:"LanguageOfPerformance,omitempty"`
	LanguageOfDubbing              string                          `xml:"LanguageOfDubbing,omitempty" json:"LanguageOfDubbing,omitempty"`
	SubTitleLanguage               string                          `xml:"SubTitleLanguage,omitempty" json:"SubTitleLanguage,omitempty"`
	Duration                       string                          `xml:"Duration,omitempty" json:"Duration,omitempty"`
	RightsAgreementId              *RightsAgreementId              `xml:"RightsAgreementId,omitempty" json:"RightsAgreementId,omitempty"`
	PLine                          *PLine                          `xml:"PLine,omitempty" json:"PLine,omitempty"`
	CLine                          *CLine                          `xml:"CLine,omitempty" json:"CLine,omitempty"`
	ArtistProfilePage              *WebPage                        `xml:"ArtistProfilePage,omitempty" json:"ArtistProfilePage,omitempty"`
	GlobalReleaseDate              *EventDate                      `xml:"GlobalReleaseDate,omitempty" json:"GlobalReleaseDate,omitempty"`
	GlobalOriginalReleaseDate      *EventDate                      `xml:"GlobalOriginalReleaseDate,omitempty" json:"GlobalOriginalReleaseDate,omitempty"`
	LanguageAndScriptCode          string                          `xml:"LanguageAndScriptCode,attr,omitempty" json:"LanguageAndScriptCode,omitempty"`
	IsMainRelease                  bool                            `xml:"IsMainRelease,attr,omitempty" json:"IsMainRelease,omitempty"`
}

type WebPage struct {
	PartyId   *PartyId   `xml:"PartyId,omitempty" json:"PartyId,omitempty"`
	ReleaseId *ReleaseId `xml:"ReleaseId,omitempty" json:"ReleaseId,omitempty"`
	PageName  *Name      `xml:"PageName,omitempty" json:"PageName,omitempty"`
	URL       string     `xml:"URL,omitempty" json:"URL,omitempty"`
	UserName  string     `xml:"UserName,omitempty" json:"UserName,omitempty"`
	Password  string     `xml:"Password,omitempty" json:"Password,omitempty"`
}

type ReleaseDetailsByTerritory struct {
	TerritoryCode                   *TerritoryCode               `xml:"TerritoryCode,omitempty" json:"TerritoryCode,omitempty"`
	ExcludedTerritoryCode           *TerritoryCode               `xml:"ExcludedTerritoryCode,omitempty" json:"ExcludedTerritoryCode,omitempty"`
	DisplayArtistName               *Name                        `xml:"DisplayArtistName,omitempty" json:"DisplayArtistName,omitempty"`
	LabelName                       []*LabelName                 `xml:"LabelName,omitempty" json:"LabelName,omitempty"`
	RightsAgreementId               *RightsAgreementId           `xml:"RightsAgreementId,omitempty" json:"RightsAgreementId,omitempty"`
	Title                           []*Title                     `xml:"Title,omitempty" json:"Title,omitempty"`
	DisplayArtist                   []*Artist                    `xml:"DisplayArtist,omitempty" json:"DisplayArtist,omitempty"`
	IsMultiArtistCompilation        bool                         `xml:"IsMultiArtistCompilation,omitempty" json:"IsMultiArtistCompilation,omitempty"`
	AdministratingRecordCompany     *AdministratingRecordCompany `xml:"AdministratingRecordCompany,omitempty" json:"AdministratingRecordCompany,omitempty"`
	ReleaseType                     *ReleaseType                 `xml:"ReleaseType,omitempty" json:"ReleaseType,omitempty"`
	RelatedRelease                  *RelatedRelease              `xml:"RelatedRelease,omitempty" json:"RelatedRelease,omitempty"`
	ParentalWarningType             *ParentalWarningType         `xml:"ParentalWarningType,omitempty" json:"ParentalWarningType,omitempty"`
	AvRating                        *AvRating                    `xml:"AvRating,omitempty" json:"AvRating,omitempty"`
	MarketingComment                *Comment                     `xml:"MarketingComment,omitempty" json:"MarketingComment,omitempty"`
	ResourceGroup                   *ResourceGroup               `xml:"ResourceGroup,omitempty" json:"ResourceGroup,omitempty"`
	Genre                           *Genre                       `xml:"Genre,omitempty" json:"Genre,omitempty"`
	PLine                           *PLine                       `xml:"PLine,omitempty" json:"PLine,omitempty"`
	CLine                           *CLine                       `xml:"CLine,omitempty" json:"CLine,omitempty"`
	ReleaseDate                     *EventDate                   `xml:"ReleaseDate,omitempty" json:"ReleaseDate,omitempty"`
	OriginalReleaseDate             *EventDate                   `xml:"OriginalReleaseDate,omitempty" json:"OriginalReleaseDate,omitempty"`
	OriginalDigitalReleaseDate      *EventDate                   `xml:"OriginalDigitalReleaseDate,omitempty" json:"OriginalDigitalReleaseDate,omitempty"`
	FileAvailabilityDescription     *Description                 `xml:"FileAvailabilityDescription,omitempty" json:"FileAvailabilityDescription,omitempty"`
	File                            *File                        `xml:"File,omitempty" json:"File,omitempty"`
	Keywords                        *Keywords                    `xml:"Keywords,omitempty" json:"Keywords,omitempty"`
	Synopsis                        *Synopsis                    `xml:"Synopsis,omitempty" json:"Synopsis,omitempty"`
	Character                       *Character                   `xml:"Character,omitempty" json:"Character,omitempty"`
	NumberOfUnitsPerPhysicalRelease int                          `xml:"NumberOfUnitsPerPhysicalRelease,omitempty" json:"NumberOfUnitsPerPhysicalRelease,omitempty"`
	DisplayConductor                []*Artist                    `xml:"DisplayConductor,omitempty" json:"DisplayConductor,omitempty"`
	LanguageAndScriptCode           string                       `xml:"LanguageAndScriptCode,attr,omitempty" json:"LanguageAndScriptCode,omitempty"`
}

type ResourceGroup struct {
	Title                              []*Title                            `xml:"Title,omitempty" json:"Title,omitempty"`
	SequenceNumber                     int                                 `xml:"SequenceNumber,omitempty" json:"SequenceNumber,omitempty"`
	DisplayArtist                      []*Artist                           `xml:"DisplayArtist,omitempty" json:"DisplayArtist,omitempty"`
	DisplayConductor                   []*Artist                           `xml:"DisplayConductor,omitempty" json:"DisplayConductor,omitempty"`
	DisplayComposer                    []*Artist                           `xml:"DisplayComposer,omitempty" json:"DisplayComposer,omitempty"`
	ResourceContributor                []*DetailedResourceContributor      `xml:"ResourceContributor,omitempty" json:"ResourceContributor,omitempty"`
	IndirectResourceContributor        []*IndirectResourceContributor      `xml:"IndirectResourceContributor,omitempty" json:"IndirectResourceContributor,omitempty"`
	CarrierType                        *CarrierType                        `xml:"CarrierType,omitempty" json:"CarrierType,omitempty"`
	ResourceGroup                      []*ResourceGroup                    `xml:"ResourceGroup,omitempty" json:"ResourceGroup,omitempty"`
	ResourceGroupContentItem           []*ExtendedResourceGroupContentItem `xml:"ResourceGroupContentItem,omitempty" json:"ResourceGroupContentItem,omitempty"`
	ResourceGroupResourceReferenceList *ResourceGroupResourceReferenceList `xml:"ResourceGroupResourceReferenceList,omitempty" json:"ResourceGroupResourceReferenceList,omitempty"`
	ResourceGroupReleaseReference      string                              `xml:"ResourceGroupReleaseReference,omitempty" json:"ResourceGroupReleaseReference,omitempty"`
	ReleaseId                          *ReleaseId                          `xml:"ReleaseId,omitempty" json:"ReleaseId,omitempty"`
	LanguageAndScriptCode              string                              `xml:"LanguageAndScriptCode,attr,omitempty" json:"LanguageAndScriptCode,omitempty"`
}

type ResourceGroupResourceReferenceList struct {
	ResourceGroupResourceReference []string `xml:"ResourceGroupResourceReference,omitempty" json:"ResourceGroupResourceReference,omitempty"`
}

type ExtendedResourceGroupContentItem struct {
	SequenceNumber                           int                             `xml:"SequenceNumber,omitempty" json:"SequenceNumber,omitempty"`
	SequenceSubNumber                        int                             `xml:"SequenceSubNumber,omitempty" json:"SequenceSubNumber,omitempty"`
	ResourceType                             *ResourceType                   `xml:"ResourceType,omitempty" json:"ResourceType,omitempty"`
	ReleaseResourceReference                 *ReleaseResourceReference       `xml:"ReleaseResourceReference,omitempty" json:"ReleaseResourceReference,omitempty"`
	LinkedReleaseResourceReference           *LinkedReleaseResourceReference `xml:"LinkedReleaseResourceReference,omitempty" json:"LinkedReleaseResourceReference,omitempty"`
	ResourceGroupContentItemReleaseReference string                          `xml:"ResourceGroupContentItemReleaseReference,omitempty" json:"ResourceGroupContentItemReleaseReference,omitempty"`
	ReleaseId                                *ReleaseId                      `xml:"ReleaseId,omitempty" json:"ReleaseId,omitempty"`
	Duration                                 string                          `xml:"Duration,omitempty" json:"Duration,omitempty"`
	IsHiddenResource                         bool                            `xml:"IsHiddenResource,omitempty" json:"IsHiddenResource,omitempty"`
	IsBonusResource                          bool                            `xml:"IsBonusResource,omitempty" json:"IsBonusResource,omitempty"`
	IsInstantGratificationResource           bool                            `xml:"IsInstantGratificationResource,omitempty" json:"IsInstantGratificationResource,omitempty"`
	IsPreOrderIncentiveResource              bool                            `xml:"IsPreOrderIncentiveResource,omitempty" json:"IsPreOrderIncentiveResource,omitempty"`
}

type LinkedReleaseResourceReference struct {
	Value                 string `xml:",chardata"`
	LinkDescription       string `xml:"LinkDescription,attr,omitempty" json:"LinkDescription,omitempty"`
	LanguageAndScriptCode string `xml:"LanguageAndScriptCode,attr,omitempty" json:"LanguageAndScriptCode,omitempty"`
}

type ResourceType struct {
	Value            string `xml:",chardata"`
	Namespace        string `xml:"Namespace,attr,omitempty" json:"Namespace,omitempty"`
	UserDefinedValue string `xml:"UserDefinedValue,attr,omitempty" json:"UserDefinedValue,omitempty"`
}

type RelatedRelease struct {
	ReleaseId                        *ReleaseId                          `xml:"ReleaseId,omitempty" json:"ReleaseId,omitempty"`
	ReferenceTitle                   *ReferenceTitle                     `xml:"ReferenceTitle,omitempty" json:"ReferenceTitle,omitempty"`
	ReleaseSummaryDetailsByTerritory []*ReleaseSummaryDetailsByTerritory `xml:"ReleaseSummaryDetailsByTerritory,omitempty" json:"ReleaseSummaryDetailsByTerritory,omitempty"`
	RightsAgreementId                *RightsAgreementId                  `xml:"RightsAgreementId,omitempty" json:"RightsAgreementId,omitempty"`
	ReleaseRelationshipType          *ReleaseRelationshipType            `xml:"ReleaseRelationshipType,omitempty" json:"ReleaseRelationshipType,omitempty"`
	ReleaseDate                      *EventDate                          `xml:"ReleaseDate,omitempty" json:"ReleaseDate,omitempty"`
	OriginalReleaseDate              *EventDate                          `xml:"OriginalReleaseDate,omitempty" json:"OriginalReleaseDate,omitempty"`
	LanguageAndScriptCode            string                              `xml:"LanguageAndScriptCode,attr,omitempty" json:"LanguageAndScriptCode,omitempty"`
}

type ReleaseRelationshipType struct {
	Value            string `xml:",chardata"`
	Namespace        string `xml:"Namespace,attr,omitempty" json:"Namespace,omitempty"`
	UserDefinedValue string `xml:"UserDefinedValue,attr,omitempty" json:"UserDefinedValue,omitempty"`
}

type ReleaseSummaryDetailsByTerritory struct {
	TerritoryCode         *TerritoryCode     `xml:"TerritoryCode,omitempty" json:"TerritoryCode,omitempty"`
	ExcludedTerritoryCode *TerritoryCode     `xml:"ExcludedTerritoryCode,omitempty" json:"ExcludedTerritoryCode,omitempty"`
	DisplayArtistName     *Name              `xml:"DisplayArtistName,omitempty" json:"DisplayArtistName,omitempty"`
	LabelName             []*LabelName       `xml:"LabelName,omitempty" json:"LabelName,omitempty"`
	RightsAgreementId     *RightsAgreementId `xml:"RightsAgreementId,omitempty" json:"RightsAgreementId,omitempty"`
	LanguageAndScriptCode string             `xml:"LanguageAndScriptCode,attr,omitempty" json:"LanguageAndScriptCode,omitempty"`
}

type ReleaseType struct {
	Value            string `xml:",chardata"`
	Namespace        string `xml:"Namespace,attr,omitempty" json:"Namespace,omitempty"`
	UserDefinedValue string `xml:"UserDefinedValue,attr,omitempty" json:"UserDefinedValue,omitempty"`
}

type ReleaseCollectionReferenceList struct {
	NumberOfCollections        int                           `xml:"NumberOfCollections,omitempty" json:"NumberOfCollections,omitempty"`
	ReleaseCollectionReference []*ReleaseCollectionReference `xml:"ReleaseCollectionReference,omitempty" json:"ReleaseCollectionReference,omitempty"`
}

type ReleaseCollectionReference struct {
	Value               string `xml:",chardata"`
	ReleaseResourceType string `xml:"ReleaseResourceType,attr,omitempty" json:"ReleaseResourceType,omitempty"`
}

type ResourceOmissionReason struct {
	Value            string `xml:",chardata"`
	Namespace        string `xml:"Namespace,attr,omitempty" json:"Namespace,omitempty"`
	UserDefinedValue string `xml:"UserDefinedValue,attr,omitempty" json:"UserDefinedValue,omitempty"`
}

type ReleaseResourceReferenceList struct {
	ReleaseResourceReference []*ReleaseResourceReference `xml:"ReleaseResourceReference,omitempty" json:"ReleaseResourceReference,omitempty"`
}

type ReleaseResourceReference struct {
	Value               string `xml:",chardata"`
	ReleaseResourceType string `xml:"ReleaseResourceType,attr,omitempty" json:"ReleaseResourceType,omitempty"`
}

type SalesReportingProxyReleaseId struct {
	ReleaseId  *ReleaseId  `xml:"ReleaseId,omitempty" json:"ReleaseId,omitempty"`
	Reason     *Reason     `xml:"Reason,omitempty" json:"Reason,omitempty"`
	ReasonType *ReasonType `xml:"ReasonType,omitempty" json:"ReasonType,omitempty"`
}

type ReasonType struct {
	Value            string `xml:",chardata"`
	Namespace        string `xml:"Namespace,attr,omitempty" json:"Namespace,omitempty"`
	UserDefinedValue string `xml:"UserDefinedValue,attr,omitempty" json:"UserDefinedValue,omitempty"`
}

type ExternalResourceLink struct {
	URL                          string                        `xml:"URL,omitempty" json:"URL,omitempty"`
	ValidityPeriod               *Period                       `xml:"ValidityPeriod,omitempty" json:"ValidityPeriod,omitempty"`
	ExternalLink                 string                        `xml:"ExternalLink,omitempty" json:"ExternalLink,omitempty"`
	ExternallyLinkedResourceType *ExternallyLinkedResourceType `xml:"ExternallyLinkedResourceType,omitempty" json:"ExternallyLinkedResourceType,omitempty"`
	FileFormat                   string                        `xml:"FileFormat,omitempty" json:"FileFormat,omitempty"`
}

type ExternallyLinkedResourceType struct {
	Value            string `xml:",chardata"`
	Namespace        string `xml:"Namespace,attr,omitempty" json:"Namespace,omitempty"`
	UserDefinedValue string `xml:"UserDefinedValue,attr,omitempty" json:"UserDefinedValue,omitempty"`
}

type CollectionList struct {
	Collection            []*Collection `xml:"Collection,omitempty" json:"Collection,omitempty"`
	LanguageAndScriptCode string        `xml:"LanguageAndScriptCode,attr,omitempty" json:"LanguageAndScriptCode,omitempty"`
}

type Collection struct {
	CollectionId                      *CollectionId                      `xml:"CollectionId,omitempty" json:"CollectionId,omitempty"`
	CollectionType                    *CollectionType                    `xml:"CollectionType,omitempty" json:"CollectionType,omitempty"`
	CollectionReference               string                             `xml:"CollectionReference,omitempty" json:"CollectionReference,omitempty"`
	EquivalentReleaseReference        string                             `xml:"EquivalentReleaseReference,omitempty" json:"EquivalentReleaseReference,omitempty"`
	Title                             []*Title                           `xml:"Title,omitempty" json:"Title,omitempty"`
	SequenceNumber                    int                                `xml:"SequenceNumber,omitempty" json:"SequenceNumber,omitempty"`
	Contributor                       []*DetailedResourceContributor     `xml:"Contributor,omitempty" json:"Contributor,omitempty"`
	Character                         *Character                         `xml:"Character,omitempty" json:"Character,omitempty"`
	CollectionCollectionReferenceList *CollectionCollectionReferenceList `xml:"CollectionCollectionReferenceList,omitempty" json:"CollectionCollectionReferenceList,omitempty"`
	IsComplete                        bool                               `xml:"IsComplete,omitempty" json:"IsComplete,omitempty"`
	Duration                          string                             `xml:"Duration,omitempty" json:"Duration,omitempty"`
	DurationOfMusicalContent          string                             `xml:"DurationOfMusicalContent,omitempty" json:"DurationOfMusicalContent,omitempty"`
	CreationDate                      *EventDate                         `xml:"CreationDate,omitempty" json:"CreationDate,omitempty"`
	ReleaseDate                       *EventDate                         `xml:"ReleaseDate,omitempty" json:"ReleaseDate,omitempty"`
	OriginalReleaseDate               *EventDate                         `xml:"OriginalReleaseDate,omitempty" json:"OriginalReleaseDate,omitempty"`
	OriginalLanguage                  string                             `xml:"OriginalLanguage,omitempty" json:"OriginalLanguage,omitempty"`
	CollectionDetailsByTerritory      []*CollectionDetailsByTerritory    `xml:"CollectionDetailsByTerritory,omitempty" json:"CollectionDetailsByTerritory,omitempty"`
	CollectionResourceReferenceList   *CollectionResourceReferenceList   `xml:"CollectionResourceReferenceList,omitempty" json:"CollectionResourceReferenceList,omitempty"`
	CollectionWorkReferenceList       *CollectionWorkReferenceList       `xml:"CollectionWorkReferenceList,omitempty" json:"CollectionWorkReferenceList,omitempty"`
	RepresentativeImageReference      string                             `xml:"RepresentativeImageReference,omitempty" json:"RepresentativeImageReference,omitempty"`
	PLine                             *PLine                             `xml:"PLine,omitempty" json:"PLine,omitempty"`
	CLine                             *CLine                             `xml:"CLine,omitempty" json:"CLine,omitempty"`
	LanguageAndScriptCode             string                             `xml:"LanguageAndScriptCode,attr,omitempty" json:"LanguageAndScriptCode,omitempty"`
}

type CollectionWorkReferenceList struct {
	CollectionWorkReference []*CollectionWorkReference `xml:"CollectionWorkReference,omitempty" json:"CollectionWorkReference,omitempty"`
}

type CollectionWorkReference struct {
	CollectionWorkReference string `xml:"CollectionWorkReference,omitempty" json:"CollectionWorkReference,omitempty"`
	Duration                string `xml:"Duration,omitempty" json:"Duration,omitempty"`
}

type CollectionResourceReferenceList struct {
	CollectionResourceReference []*CollectionResourceReference `xml:"CollectionResourceReference,omitempty" json:"CollectionResourceReference,omitempty"`
}

type CollectionResourceReference struct {
	SequenceNumber              int    `xml:"SequenceNumber,omitempty" json:"SequenceNumber,omitempty"`
	CollectionResourceReference string `xml:"CollectionResourceReference,omitempty" json:"CollectionResourceReference,omitempty"`
	Duration                    string `xml:"Duration,omitempty" json:"Duration,omitempty"`
}

type CollectionDetailsByTerritory struct {
	TerritoryCode         *TerritoryCode                 `xml:"TerritoryCode,omitempty" json:"TerritoryCode,omitempty"`
	ExcludedTerritoryCode *TerritoryCode                 `xml:"ExcludedTerritoryCode,omitempty" json:"ExcludedTerritoryCode,omitempty"`
	Title                 []*Title                       `xml:"Title,omitempty" json:"Title,omitempty"`
	Contributor           []*DetailedResourceContributor `xml:"Contributor,omitempty" json:"Contributor,omitempty"`
	IsComplete            bool                           `xml:"IsComplete,omitempty" json:"IsComplete,omitempty"`
	Character             *Character                     `xml:"Character,omitempty" json:"Character,omitempty"`
}

type CollectionCollectionReferenceList struct {
	NumberOfCollections           int                              `xml:"NumberOfCollections,omitempty" json:"NumberOfCollections,omitempty"`
	CollectionCollectionReference []*CollectionCollectionReference `xml:"CollectionCollectionReference,omitempty" json:"CollectionCollectionReference,omitempty"`
}

type CollectionCollectionReference struct {
	SequenceNumber                int    `xml:"SequenceNumber,omitempty" json:"SequenceNumber,omitempty"`
	CollectionCollectionReference string `xml:"CollectionCollectionReference,omitempty" json:"CollectionCollectionReference,omitempty"`
	StartTime                     string `xml:"StartTime,omitempty" json:"StartTime,omitempty"`
	Duration                      string `xml:"Duration,omitempty" json:"Duration,omitempty"`
	EndTime                       string `xml:"EndTime,omitempty" json:"EndTime,omitempty"`
	InclusionDate                 string `xml:"InclusionDate,omitempty" json:"InclusionDate,omitempty"`
}

type CollectionType struct {
	Value            string `xml:",chardata"`
	Namespace        string `xml:"Namespace,attr,omitempty" json:"Namespace,omitempty"`
	UserDefinedValue string `xml:"UserDefinedValue,attr,omitempty" json:"UserDefinedValue,omitempty"`
}

type CollectionId struct {
	GRid          string         `xml:"GRid,omitempty" json:"GRid,omitempty"`
	ISRC          string         `xml:"ISRC,omitempty" json:"ISRC,omitempty"`
	ISAN          string         `xml:"ISAN,omitempty" json:"ISAN,omitempty"`
	VISAN         string         `xml:"VISAN,omitempty" json:"VISAN,omitempty"`
	ICPN          *ICPN          `xml:"ICPN,omitempty" json:"ICPN,omitempty"`
	CatalogNumber *CatalogNumber `xml:"CatalogNumber,omitempty" json:"CatalogNumber,omitempty"`
	ProprietaryId *ProprietaryId `xml:"ProprietaryId,omitempty" json:"ProprietaryId,omitempty"`
	IsReplaced    bool           `xml:"IsReplaced,attr,omitempty" json:"IsReplaced,omitempty"`
}

type ResourceList struct {
	SoundRecording        []*SoundRecording      `xml:"SoundRecording,omitempty" json:"SoundRecording,omitempty"`
	MIDI                  []*MIDI                `xml:"MIDI,omitempty" json:"MIDI,omitempty"`
	Video                 []*Video               `xml:"Video,omitempty" json:"Video,omitempty"`
	Image                 []*Image               `xml:"Image,omitempty" json:"Image,omitempty"`
	Text                  []*Text                `xml:"Text,omitempty" json:"Text,omitempty"`
	SheetMusic            []*SheetMusic          `xml:"SheetMusic,omitempty" json:"SheetMusic,omitempty"`
	Software              []*Software            `xml:"Software,omitempty" json:"Software,omitempty"`
	UserDefinedResource   []*UserDefinedResource `xml:"UserDefinedResource,omitempty" json:"UserDefinedResource,omitempty"`
	LanguageAndScriptCode string                 `xml:"LanguageAndScriptCode,attr,omitempty" json:"LanguageAndScriptCode,omitempty"`
}

type UserDefinedResource struct {
	UserDefinedResourceType                *UserDefinedResourceType                 `xml:"UserDefinedResourceType,omitempty" json:"UserDefinedResourceType,omitempty"`
	IsArtistRelated                        bool                                     `xml:"IsArtistRelated,omitempty" json:"IsArtistRelated,omitempty"`
	UserDefinedResourceId                  *ResourceProprietaryId                   `xml:"UserDefinedResourceId,omitempty" json:"UserDefinedResourceId,omitempty"`
	IndirectUserDefinedResourceId          *MusicalWorkId                           `xml:"IndirectUserDefinedResourceId,omitempty" json:"IndirectUserDefinedResourceId,omitempty"`
	ResourceReference                      string                                   `xml:"ResourceReference,omitempty" json:"ResourceReference,omitempty"`
	ResourceMusicalWorkReferenceList       *ResourceMusicalWorkReferenceList        `xml:"ResourceMusicalWorkReferenceList,omitempty" json:"ResourceMusicalWorkReferenceList,omitempty"`
	ResourceContainedResourceReferenceList *ResourceContainedResourceReferenceList  `xml:"ResourceContainedResourceReferenceList,omitempty" json:"ResourceContainedResourceReferenceList,omitempty"`
	Title                                  []*Title                                 `xml:"Title,omitempty" json:"Title,omitempty"`
	UserDefinedValue                       *UserDefinedValue                        `xml:"UserDefinedValue,omitempty" json:"UserDefinedValue,omitempty"`
	CreationDate                           *EventDate                               `xml:"CreationDate,omitempty" json:"CreationDate,omitempty"`
	UserDefinedResourceDetailsByTerritory  []*UserDefinedResourceDetailsByTerritory `xml:"UserDefinedResourceDetailsByTerritory,omitempty" json:"UserDefinedResourceDetailsByTerritory,omitempty"`
	IsUpdated                              bool                                     `xml:"IsUpdated,attr,omitempty" json:"IsUpdated,omitempty"`
	LanguageAndScriptCode                  string                                   `xml:"LanguageAndScriptCode,attr,omitempty" json:"LanguageAndScriptCode,omitempty"`
}

type UserDefinedResourceDetailsByTerritory struct {
	TerritoryCode                       *TerritoryCode                       `xml:"TerritoryCode,omitempty" json:"TerritoryCode,omitempty"`
	ExcludedTerritoryCode               *TerritoryCode                       `xml:"ExcludedTerritoryCode,omitempty" json:"ExcludedTerritoryCode,omitempty"`
	Title                               []*Title                             `xml:"Title,omitempty" json:"Title,omitempty"`
	ResourceContributor                 []*DetailedResourceContributor       `xml:"ResourceContributor,omitempty" json:"ResourceContributor,omitempty"`
	IndirectResourceContributor         []*IndirectResourceContributor       `xml:"IndirectResourceContributor,omitempty" json:"IndirectResourceContributor,omitempty"`
	DisplayArtistName                   *Name                                `xml:"DisplayArtistName,omitempty" json:"DisplayArtistName,omitempty"`
	UserDefinedValue                    *UserDefinedValue                    `xml:"UserDefinedValue,omitempty" json:"UserDefinedValue,omitempty"`
	PLine                               *PLine                               `xml:"PLine,omitempty" json:"PLine,omitempty"`
	CLine                               *CLine                               `xml:"CLine,omitempty" json:"CLine,omitempty"`
	ResourceReleaseDate                 *EventDate                           `xml:"ResourceReleaseDate,omitempty" json:"ResourceReleaseDate,omitempty"`
	OriginalResourceReleaseDate         *EventDate                           `xml:"OriginalResourceReleaseDate,omitempty" json:"OriginalResourceReleaseDate,omitempty"`
	FulfillmentDate                     *FulfillmentDate                     `xml:"FulfillmentDate,omitempty" json:"FulfillmentDate,omitempty"`
	Keywords                            *Keywords                            `xml:"Keywords,omitempty" json:"Keywords,omitempty"`
	Synopsis                            *Synopsis                            `xml:"Synopsis,omitempty" json:"Synopsis,omitempty"`
	Genre                               *Genre                               `xml:"Genre,omitempty" json:"Genre,omitempty"`
	ParentalWarningType                 *ParentalWarningType                 `xml:"ParentalWarningType,omitempty" json:"ParentalWarningType,omitempty"`
	TechnicalUserDefinedResourceDetails *TechnicalUserDefinedResourceDetails `xml:"TechnicalUserDefinedResourceDetails,omitempty" json:"TechnicalUserDefinedResourceDetails,omitempty"`
	LanguageAndScriptCode               string                               `xml:"LanguageAndScriptCode,attr,omitempty" json:"LanguageAndScriptCode,omitempty"`
}

type TechnicalUserDefinedResourceDetails struct {
	TechnicalResourceDetailsReference string            `xml:"TechnicalResourceDetailsReference,omitempty" json:"TechnicalResourceDetailsReference,omitempty"`
	UserDefinedValue                  *UserDefinedValue `xml:"UserDefinedValue,omitempty" json:"UserDefinedValue,omitempty"`
	IsPreview                         bool              `xml:"IsPreview,omitempty" json:"IsPreview,omitempty"`
	PreviewDetails                    *PreviewDetails   `xml:"PreviewDetails,omitempty" json:"PreviewDetails,omitempty"`
	FulfillmentDate                   *FulfillmentDate  `xml:"FulfillmentDate,omitempty" json:"FulfillmentDate,omitempty"`
	ConsumerFulfillmentDate           *FulfillmentDate  `xml:"ConsumerFulfillmentDate,omitempty" json:"ConsumerFulfillmentDate,omitempty"`
	FileAvailabilityDescription       *Description      `xml:"FileAvailabilityDescription,omitempty" json:"FileAvailabilityDescription,omitempty"`
	File                              *File             `xml:"File,omitempty" json:"File,omitempty"`
	Fingerprint                       *Fingerprint      `xml:"Fingerprint,omitempty" json:"Fingerprint,omitempty"`
	LanguageAndScriptCode             string            `xml:"LanguageAndScriptCode,attr,omitempty" json:"LanguageAndScriptCode,omitempty"`
}

type UserDefinedValue struct {
	Value                 string `xml:",chardata"`
	Namespace             string `xml:"Namespace,attr,omitempty" json:"Namespace,omitempty"`
	Description           string `xml:"Description,attr,omitempty" json:"Description,omitempty"`
	LanguageAndScriptCode string `xml:"LanguageAndScriptCode,attr,omitempty" json:"LanguageAndScriptCode,omitempty"`
}

type UserDefinedResourceType struct {
	Value     string `xml:",chardata"`
	Namespace string `xml:"Namespace,attr,omitempty" json:"Namespace,omitempty"`
}

type Software struct {
	SoftwareType                           *SoftwareType                           `xml:"SoftwareType,omitempty" json:"SoftwareType,omitempty"`
	IsArtistRelated                        bool                                    `xml:"IsArtistRelated,omitempty" json:"IsArtistRelated,omitempty"`
	SoftwareId                             *ResourceProprietaryId                  `xml:"SoftwareId,omitempty" json:"SoftwareId,omitempty"`
	IndirectSoftwareId                     *MusicalWorkId                          `xml:"IndirectSoftwareId,omitempty" json:"IndirectSoftwareId,omitempty"`
	ResourceReference                      string                                  `xml:"ResourceReference,omitempty" json:"ResourceReference,omitempty"`
	ResourceMusicalWorkReferenceList       *ResourceMusicalWorkReferenceList       `xml:"ResourceMusicalWorkReferenceList,omitempty" json:"ResourceMusicalWorkReferenceList,omitempty"`
	ResourceContainedResourceReferenceList *ResourceContainedResourceReferenceList `xml:"ResourceContainedResourceReferenceList,omitempty" json:"ResourceContainedResourceReferenceList,omitempty"`
	Title                                  []*Title                                `xml:"Title,omitempty" json:"Title,omitempty"`
	CreationDate                           *EventDate                              `xml:"CreationDate,omitempty" json:"CreationDate,omitempty"`
	SoftwareDetailsByTerritory             []*SoftwareDetailsByTerritory           `xml:"SoftwareDetailsByTerritory,omitempty" json:"SoftwareDetailsByTerritory,omitempty"`
	IsUpdated                              bool                                    `xml:"IsUpdated,attr,omitempty" json:"IsUpdated,omitempty"`
	LanguageAndScriptCode                  string                                  `xml:"LanguageAndScriptCode,attr,omitempty" json:"LanguageAndScriptCode,omitempty"`
}

type SoftwareDetailsByTerritory struct {
	TerritoryCode               *TerritoryCode                 `xml:"TerritoryCode,omitempty" json:"TerritoryCode,omitempty"`
	ExcludedTerritoryCode       *TerritoryCode                 `xml:"ExcludedTerritoryCode,omitempty" json:"ExcludedTerritoryCode,omitempty"`
	Title                       []*Title                       `xml:"Title,omitempty" json:"Title,omitempty"`
	ResourceContributor         []*DetailedResourceContributor `xml:"ResourceContributor,omitempty" json:"ResourceContributor,omitempty"`
	IndirectResourceContributor []*IndirectResourceContributor `xml:"IndirectResourceContributor,omitempty" json:"IndirectResourceContributor,omitempty"`
	DisplayArtistName           *Name                          `xml:"DisplayArtistName,omitempty" json:"DisplayArtistName,omitempty"`
	PLine                       *PLine                         `xml:"PLine,omitempty" json:"PLine,omitempty"`
	CLine                       *CLine                         `xml:"CLine,omitempty" json:"CLine,omitempty"`
	CourtesyLine                *CourtesyLine                  `xml:"CourtesyLine,omitempty" json:"CourtesyLine,omitempty"`
	ResourceReleaseDate         *EventDate                     `xml:"ResourceReleaseDate,omitempty" json:"ResourceReleaseDate,omitempty"`
	OriginalResourceReleaseDate *EventDate                     `xml:"OriginalResourceReleaseDate,omitempty" json:"OriginalResourceReleaseDate,omitempty"`
	FulfillmentDate             *FulfillmentDate               `xml:"FulfillmentDate,omitempty" json:"FulfillmentDate,omitempty"`
	Keywords                    *Keywords                      `xml:"Keywords,omitempty" json:"Keywords,omitempty"`
	Synopsis                    *Synopsis                      `xml:"Synopsis,omitempty" json:"Synopsis,omitempty"`
	Genre                       *Genre                         `xml:"Genre,omitempty" json:"Genre,omitempty"`
	ParentalWarningType         *ParentalWarningType           `xml:"ParentalWarningType,omitempty" json:"ParentalWarningType,omitempty"`
	TechnicalSoftwareDetails    *TechnicalSoftwareDetails      `xml:"TechnicalSoftwareDetails,omitempty" json:"TechnicalSoftwareDetails,omitempty"`
	LanguageAndScriptCode       string                         `xml:"LanguageAndScriptCode,attr,omitempty" json:"LanguageAndScriptCode,omitempty"`
}

type TechnicalSoftwareDetails struct {
	TechnicalResourceDetailsReference string               `xml:"TechnicalResourceDetailsReference,omitempty" json:"TechnicalResourceDetailsReference,omitempty"`
	DrmPlatformType                   *DrmPlatformType     `xml:"DrmPlatformType,omitempty" json:"DrmPlatformType,omitempty"`
	OperatingSystemType               *OperatingSystemType `xml:"OperatingSystemType,omitempty" json:"OperatingSystemType,omitempty"`
	IsPreview                         bool                 `xml:"IsPreview,omitempty" json:"IsPreview,omitempty"`
	PreviewDetails                    *PreviewDetails      `xml:"PreviewDetails,omitempty" json:"PreviewDetails,omitempty"`
	FulfillmentDate                   *FulfillmentDate     `xml:"FulfillmentDate,omitempty" json:"FulfillmentDate,omitempty"`
	ConsumerFulfillmentDate           *FulfillmentDate     `xml:"ConsumerFulfillmentDate,omitempty" json:"ConsumerFulfillmentDate,omitempty"`
	FileAvailabilityDescription       *Description         `xml:"FileAvailabilityDescription,omitempty" json:"FileAvailabilityDescription,omitempty"`
	File                              *File                `xml:"File,omitempty" json:"File,omitempty"`
	Fingerprint                       *Fingerprint         `xml:"Fingerprint,omitempty" json:"Fingerprint,omitempty"`
	LanguageAndScriptCode             string               `xml:"LanguageAndScriptCode,attr,omitempty" json:"LanguageAndScriptCode,omitempty"`
}

type OperatingSystemType struct {
	Value            string `xml:",chardata"`
	Version          string `xml:"Version,attr,omitempty" json:"Version,omitempty"`
	Namespace        string `xml:"Namespace,attr,omitempty" json:"Namespace,omitempty"`
	UserDefinedValue string `xml:"UserDefinedValue,attr,omitempty" json:"UserDefinedValue,omitempty"`
}

type SoftwareType struct {
	Value            string `xml:",chardata"`
	Namespace        string `xml:"Namespace,attr,omitempty" json:"Namespace,omitempty"`
	UserDefinedValue string `xml:"UserDefinedValue,attr,omitempty" json:"UserDefinedValue,omitempty"`
}

type SheetMusic struct {
	SheetMusicType                         *SheetMusicType                         `xml:"SheetMusicType,omitempty" json:"SheetMusicType,omitempty"`
	IsArtistRelated                        bool                                    `xml:"IsArtistRelated,omitempty" json:"IsArtistRelated,omitempty"`
	SheetMusicId                           *SheetMusicId                           `xml:"SheetMusicId,omitempty" json:"SheetMusicId,omitempty"`
	IndirectSheetMusicId                   *MusicalWorkId                          `xml:"IndirectSheetMusicId,omitempty" json:"IndirectSheetMusicId,omitempty"`
	ResourceReference                      string                                  `xml:"ResourceReference,omitempty" json:"ResourceReference,omitempty"`
	LanguageOfLyrics                       string                                  `xml:"LanguageOfLyrics,omitempty" json:"LanguageOfLyrics,omitempty"`
	RightsAgreementId                      *RightsAgreementId                      `xml:"RightsAgreementId,omitempty" json:"RightsAgreementId,omitempty"`
	ResourceMusicalWorkReferenceList       *ResourceMusicalWorkReferenceList       `xml:"ResourceMusicalWorkReferenceList,omitempty" json:"ResourceMusicalWorkReferenceList,omitempty"`
	ResourceContainedResourceReferenceList *ResourceContainedResourceReferenceList `xml:"ResourceContainedResourceReferenceList,omitempty" json:"ResourceContainedResourceReferenceList,omitempty"`
	ReferenceTitle                         *ReferenceTitle                         `xml:"ReferenceTitle,omitempty" json:"ReferenceTitle,omitempty"`
	CreationDate                           *EventDate                              `xml:"CreationDate,omitempty" json:"CreationDate,omitempty"`
	SheetMusicDetailsByTerritory           []*SheetMusicDetailsByTerritory         `xml:"SheetMusicDetailsByTerritory,omitempty" json:"SheetMusicDetailsByTerritory,omitempty"`
	IsUpdated                              bool                                    `xml:"IsUpdated,attr,omitempty" json:"IsUpdated,omitempty"`
	LanguageAndScriptCode                  string                                  `xml:"LanguageAndScriptCode,attr,omitempty" json:"LanguageAndScriptCode,omitempty"`
}

type SheetMusicDetailsByTerritory struct {
	TerritoryCode               *TerritoryCode                 `xml:"TerritoryCode,omitempty" json:"TerritoryCode,omitempty"`
	ExcludedTerritoryCode       *TerritoryCode                 `xml:"ExcludedTerritoryCode,omitempty" json:"ExcludedTerritoryCode,omitempty"`
	Title                       []*Title                       `xml:"Title,omitempty" json:"Title,omitempty"`
	ResourceContributor         []*DetailedResourceContributor `xml:"ResourceContributor,omitempty" json:"ResourceContributor,omitempty"`
	IndirectResourceContributor []*IndirectResourceContributor `xml:"IndirectResourceContributor,omitempty" json:"IndirectResourceContributor,omitempty"`
	DisplayArtistName           *Name                          `xml:"DisplayArtistName,omitempty" json:"DisplayArtistName,omitempty"`
	CLine                       *CLine                         `xml:"CLine,omitempty" json:"CLine,omitempty"`
	CourtesyLine                *CourtesyLine                  `xml:"CourtesyLine,omitempty" json:"CourtesyLine,omitempty"`
	ResourceReleaseDate         *EventDate                     `xml:"ResourceReleaseDate,omitempty" json:"ResourceReleaseDate,omitempty"`
	OriginalResourceReleaseDate *EventDate                     `xml:"OriginalResourceReleaseDate,omitempty" json:"OriginalResourceReleaseDate,omitempty"`
	FulfillmentDate             *FulfillmentDate               `xml:"FulfillmentDate,omitempty" json:"FulfillmentDate,omitempty"`
	Genre                       *Genre                         `xml:"Genre,omitempty" json:"Genre,omitempty"`
	ParentalWarningType         *ParentalWarningType           `xml:"ParentalWarningType,omitempty" json:"ParentalWarningType,omitempty"`
	TechnicalSheetMusicDetails  *TechnicalSheetMusicDetails    `xml:"TechnicalSheetMusicDetails,omitempty" json:"TechnicalSheetMusicDetails,omitempty"`
	LanguageAndScriptCode       string                         `xml:"LanguageAndScriptCode,attr,omitempty" json:"LanguageAndScriptCode,omitempty"`
}

type TechnicalSheetMusicDetails struct {
	TechnicalResourceDetailsReference string               `xml:"TechnicalResourceDetailsReference,omitempty" json:"TechnicalResourceDetailsReference,omitempty"`
	DrmPlatformType                   *DrmPlatformType     `xml:"DrmPlatformType,omitempty" json:"DrmPlatformType,omitempty"`
	ContainerFormat                   *ContainerFormat     `xml:"ContainerFormat,omitempty" json:"ContainerFormat,omitempty"`
	SheetMusicCodecType               *SheetMusicCodecType `xml:"SheetMusicCodecType,omitempty" json:"SheetMusicCodecType,omitempty"`
	IsPreview                         bool                 `xml:"IsPreview,omitempty" json:"IsPreview,omitempty"`
	PreviewDetails                    *PreviewDetails      `xml:"PreviewDetails,omitempty" json:"PreviewDetails,omitempty"`
	FulfillmentDate                   *FulfillmentDate     `xml:"FulfillmentDate,omitempty" json:"FulfillmentDate,omitempty"`
	ConsumerFulfillmentDate           *FulfillmentDate     `xml:"ConsumerFulfillmentDate,omitempty" json:"ConsumerFulfillmentDate,omitempty"`
	FileAvailabilityDescription       *Description         `xml:"FileAvailabilityDescription,omitempty" json:"FileAvailabilityDescription,omitempty"`
	File                              *File                `xml:"File,omitempty" json:"File,omitempty"`
	Fingerprint                       *Fingerprint         `xml:"Fingerprint,omitempty" json:"Fingerprint,omitempty"`
	LanguageAndScriptCode             string               `xml:"LanguageAndScriptCode,attr,omitempty" json:"LanguageAndScriptCode,omitempty"`
}

type SheetMusicCodecType struct {
	Value            string `xml:",chardata"`
	Version          string `xml:"Version,attr,omitempty" json:"Version,omitempty"`
	Namespace        string `xml:"Namespace,attr,omitempty" json:"Namespace,omitempty"`
	UserDefinedValue string `xml:"UserDefinedValue,attr,omitempty" json:"UserDefinedValue,omitempty"`
}

type SheetMusicId struct {
	ISMN          string         `xml:"ISMN,omitempty" json:"ISMN,omitempty"`
	ProprietaryId *ProprietaryId `xml:"ProprietaryId,omitempty" json:"ProprietaryId,omitempty"`
	IsReplaced    bool           `xml:"IsReplaced,attr,omitempty" json:"IsReplaced,omitempty"`
}

type SheetMusicType struct {
	Value            string `xml:",chardata"`
	Namespace        string `xml:"Namespace,attr,omitempty" json:"Namespace,omitempty"`
	UserDefinedValue string `xml:"UserDefinedValue,attr,omitempty" json:"UserDefinedValue,omitempty"`
}

type Text struct {
	TextType                               *TextType                               `xml:"TextType,omitempty" json:"TextType,omitempty"`
	IsArtistRelated                        bool                                    `xml:"IsArtistRelated,omitempty" json:"IsArtistRelated,omitempty"`
	TextId                                 *TextId                                 `xml:"TextId,omitempty" json:"TextId,omitempty"`
	IndirectTextId                         *MusicalWorkId                          `xml:"IndirectTextId,omitempty" json:"IndirectTextId,omitempty"`
	ResourceReference                      string                                  `xml:"ResourceReference,omitempty" json:"ResourceReference,omitempty"`
	ResourceMusicalWorkReferenceList       *ResourceMusicalWorkReferenceList       `xml:"ResourceMusicalWorkReferenceList,omitempty" json:"ResourceMusicalWorkReferenceList,omitempty"`
	ResourceContainedResourceReferenceList *ResourceContainedResourceReferenceList `xml:"ResourceContainedResourceReferenceList,omitempty" json:"ResourceContainedResourceReferenceList,omitempty"`
	Title                                  []*Title                                `xml:"Title,omitempty" json:"Title,omitempty"`
	CreationDate                           *EventDate                              `xml:"CreationDate,omitempty" json:"CreationDate,omitempty"`
	TextDetailsByTerritory                 []*TextDetailsByTerritory               `xml:"TextDetailsByTerritory,omitempty" json:"TextDetailsByTerritory,omitempty"`
	IsUpdated                              bool                                    `xml:"IsUpdated,attr,omitempty" json:"IsUpdated,omitempty"`
	LanguageAndScriptCode                  string                                  `xml:"LanguageAndScriptCode,attr,omitempty" json:"LanguageAndScriptCode,omitempty"`
}

type TextDetailsByTerritory struct {
	TerritoryCode               *TerritoryCode                 `xml:"TerritoryCode,omitempty" json:"TerritoryCode,omitempty"`
	ExcludedTerritoryCode       *TerritoryCode                 `xml:"ExcludedTerritoryCode,omitempty" json:"ExcludedTerritoryCode,omitempty"`
	Title                       []*Title                       `xml:"Title,omitempty" json:"Title,omitempty"`
	ResourceContributor         []*DetailedResourceContributor `xml:"ResourceContributor,omitempty" json:"ResourceContributor,omitempty"`
	IndirectResourceContributor []*IndirectResourceContributor `xml:"IndirectResourceContributor,omitempty" json:"IndirectResourceContributor,omitempty"`
	DisplayArtistName           *Name                          `xml:"DisplayArtistName,omitempty" json:"DisplayArtistName,omitempty"`
	CLine                       *CLine                         `xml:"CLine,omitempty" json:"CLine,omitempty"`
	CourtesyLine                *CourtesyLine                  `xml:"CourtesyLine,omitempty" json:"CourtesyLine,omitempty"`
	ResourceReleaseDate         *EventDate                     `xml:"ResourceReleaseDate,omitempty" json:"ResourceReleaseDate,omitempty"`
	OriginalResourceReleaseDate *EventDate                     `xml:"OriginalResourceReleaseDate,omitempty" json:"OriginalResourceReleaseDate,omitempty"`
	FulfillmentDate             *FulfillmentDate               `xml:"FulfillmentDate,omitempty" json:"FulfillmentDate,omitempty"`
	Keywords                    *Keywords                      `xml:"Keywords,omitempty" json:"Keywords,omitempty"`
	Synopsis                    *Synopsis                      `xml:"Synopsis,omitempty" json:"Synopsis,omitempty"`
	Genre                       *Genre                         `xml:"Genre,omitempty" json:"Genre,omitempty"`
	ParentalWarningType         *ParentalWarningType           `xml:"ParentalWarningType,omitempty" json:"ParentalWarningType,omitempty"`
	TechnicalTextDetails        *TechnicalTextDetails          `xml:"TechnicalTextDetails,omitempty" json:"TechnicalTextDetails,omitempty"`
	LanguageAndScriptCode       string                         `xml:"LanguageAndScriptCode,attr,omitempty" json:"LanguageAndScriptCode,omitempty"`
}

type TechnicalTextDetails struct {
	TechnicalResourceDetailsReference string           `xml:"TechnicalResourceDetailsReference,omitempty" json:"TechnicalResourceDetailsReference,omitempty"`
	DrmPlatformType                   *DrmPlatformType `xml:"DrmPlatformType,omitempty" json:"DrmPlatformType,omitempty"`
	ContainerFormat                   *ContainerFormat `xml:"ContainerFormat,omitempty" json:"ContainerFormat,omitempty"`
	TextCodecType                     *TextCodecType   `xml:"TextCodecType,omitempty" json:"TextCodecType,omitempty"`
	IsPreview                         bool             `xml:"IsPreview,omitempty" json:"IsPreview,omitempty"`
	PreviewDetails                    *PreviewDetails  `xml:"PreviewDetails,omitempty" json:"PreviewDetails,omitempty"`
	FulfillmentDate                   *FulfillmentDate `xml:"FulfillmentDate,omitempty" json:"FulfillmentDate,omitempty"`
	ConsumerFulfillmentDate           *FulfillmentDate `xml:"ConsumerFulfillmentDate,omitempty" json:"ConsumerFulfillmentDate,omitempty"`
	FileAvailabilityDescription       *Description     `xml:"FileAvailabilityDescription,omitempty" json:"FileAvailabilityDescription,omitempty"`
	File                              *File            `xml:"File,omitempty" json:"File,omitempty"`
	Fingerprint                       *Fingerprint     `xml:"Fingerprint,omitempty" json:"Fingerprint,omitempty"`
	LanguageAndScriptCode             string           `xml:"LanguageAndScriptCode,attr,omitempty" json:"LanguageAndScriptCode,omitempty"`
}

type TextCodecType struct {
	Value            string `xml:",chardata"`
	Version          string `xml:"Version,attr,omitempty" json:"Version,omitempty"`
	Namespace        string `xml:"Namespace,attr,omitempty" json:"Namespace,omitempty"`
	UserDefinedValue string `xml:"UserDefinedValue,attr,omitempty" json:"UserDefinedValue,omitempty"`
}

type TextId struct {
	ISBN          string         `xml:"ISBN,omitempty" json:"ISBN,omitempty"`
	ISSN          string         `xml:"ISSN,omitempty" json:"ISSN,omitempty"`
	SICI          string         `xml:"SICI,omitempty" json:"SICI,omitempty"`
	ProprietaryId *ProprietaryId `xml:"ProprietaryId,omitempty" json:"ProprietaryId,omitempty"`
	IsReplaced    bool           `xml:"IsReplaced,attr,omitempty" json:"IsReplaced,omitempty"`
}

type TextType struct {
	Value            string `xml:",chardata"`
	Namespace        string `xml:"Namespace,attr,omitempty" json:"Namespace,omitempty"`
	UserDefinedValue string `xml:"UserDefinedValue,attr,omitempty" json:"UserDefinedValue,omitempty"`
}

type Image struct {
	ImageType               *ImageType                 `xml:"ImageType,omitempty" json:"ImageType,omitempty"`
	IsArtistRelated         bool                       `xml:"IsArtistRelated,omitempty" json:"IsArtistRelated,omitempty"`
	ImageId                 *ResourceProprietaryId     `xml:"ImageId,omitempty" json:"ImageId,omitempty"`
	ResourceReference       string                     `xml:"ResourceReference,omitempty" json:"ResourceReference,omitempty"`
	Title                   []*Title                   `xml:"Title,omitempty" json:"Title,omitempty"`
	CreationDate            *EventDate                 `xml:"CreationDate,omitempty" json:"CreationDate,omitempty"`
	ImageDetailsByTerritory []*ImageDetailsByTerritory `xml:"ImageDetailsByTerritory,omitempty" json:"ImageDetailsByTerritory,omitempty"`
	IsUpdated               bool                       `xml:"IsUpdated,attr,omitempty" json:"IsUpdated,omitempty"`
	LanguageAndScriptCode   string                     `xml:"LanguageAndScriptCode,attr,omitempty" json:"LanguageAndScriptCode,omitempty"`
}

type Title struct {
	TitleType             string           `xml:"TitleType,attr,omitempty" json:"TitleType,omitempty"`
	TitleText             *TitleText       `xml:"TitleText,omitempty" json:"TitleText,omitempty"`
	SubTitle              []*TypedSubTitle `xml:"SubTitle,omitempty" json:"SubTitle,omitempty"`
	LanguageAndScriptCode string           `xml:"LanguageAndScriptCode,attr,omitempty" json:"LanguageAndScriptCode,omitempty"`
}

type TypedSubTitle struct {
	Value                 string `xml:",chardata"`
	LanguageAndScriptCode string `xml:"LanguageAndScriptCode,attr,omitempty" json:"LanguageAndScriptCode,omitempty"`
	SubTitleType          string `xml:"SubTitleType,attr,omitempty" json:"SubTitleType,omitempty"`
}

type ImageDetailsByTerritory struct {
	TerritoryCode               *TerritoryCode                 `xml:"TerritoryCode,omitempty" json:"TerritoryCode,omitempty"`
	ExcludedTerritoryCode       *TerritoryCode                 `xml:"ExcludedTerritoryCode,omitempty" json:"ExcludedTerritoryCode,omitempty"`
	Title                       []*Title                       `xml:"Title,omitempty" json:"Title,omitempty"`
	ResourceContributor         []*DetailedResourceContributor `xml:"ResourceContributor,omitempty" json:"ResourceContributor,omitempty"`
	IndirectResourceContributor []*IndirectResourceContributor `xml:"IndirectResourceContributor,omitempty" json:"IndirectResourceContributor,omitempty"`
	DisplayArtistName           *Name                          `xml:"DisplayArtistName,omitempty" json:"DisplayArtistName,omitempty"`
	CLine                       *CLine                         `xml:"CLine,omitempty" json:"CLine,omitempty"`
	Description                 *Description                   `xml:"Description,omitempty" json:"Description,omitempty"`
	CourtesyLine                *CourtesyLine                  `xml:"CourtesyLine,omitempty" json:"CourtesyLine,omitempty"`
	ResourceReleaseDate         *EventDate                     `xml:"ResourceReleaseDate,omitempty" json:"ResourceReleaseDate,omitempty"`
	OriginalResourceReleaseDate *EventDate                     `xml:"OriginalResourceReleaseDate,omitempty" json:"OriginalResourceReleaseDate,omitempty"`
	FulfillmentDate             *FulfillmentDate               `xml:"FulfillmentDate,omitempty" json:"FulfillmentDate,omitempty"`
	Keywords                    *Keywords                      `xml:"Keywords,omitempty" json:"Keywords,omitempty"`
	Synopsis                    *Synopsis                      `xml:"Synopsis,omitempty" json:"Synopsis,omitempty"`
	Genre                       *Genre                         `xml:"Genre,omitempty" json:"Genre,omitempty"`
	ParentalWarningType         *ParentalWarningType           `xml:"ParentalWarningType,omitempty" json:"ParentalWarningType,omitempty"`
	TechnicalImageDetails       *TechnicalImageDetails         `xml:"TechnicalImageDetails,omitempty" json:"TechnicalImageDetails,omitempty"`
	LanguageAndScriptCode       string                         `xml:"LanguageAndScriptCode,attr,omitempty" json:"LanguageAndScriptCode,omitempty"`
}

type TechnicalImageDetails struct {
	TechnicalResourceDetailsReference string           `xml:"TechnicalResourceDetailsReference,omitempty" json:"TechnicalResourceDetailsReference,omitempty"`
	DrmPlatformType                   *DrmPlatformType `xml:"DrmPlatformType,omitempty" json:"DrmPlatformType,omitempty"`
	ContainerFormat                   *ContainerFormat `xml:"ContainerFormat,omitempty" json:"ContainerFormat,omitempty"`
	ImageCodecType                    *ImageCodecType  `xml:"ImageCodecType,omitempty" json:"ImageCodecType,omitempty"`
	ImageHeight                       *Extent          `xml:"ImageHeight,omitempty" json:"ImageHeight,omitempty"`
	ImageWidth                        *Extent          `xml:"ImageWidth,omitempty" json:"ImageWidth,omitempty"`
	AspectRatio                       *AspectRatio     `xml:"AspectRatio,omitempty" json:"AspectRatio,omitempty"`
	ColorDepth                        int              `xml:"ColorDepth,omitempty" json:"ColorDepth,omitempty"`
	ImageResolution                   int              `xml:"ImageResolution,omitempty" json:"ImageResolution,omitempty"`
	IsPreview                         bool             `xml:"IsPreview,omitempty" json:"IsPreview,omitempty"`
	PreviewDetails                    *PreviewDetails  `xml:"PreviewDetails,omitempty" json:"PreviewDetails,omitempty"`
	FulfillmentDate                   *FulfillmentDate `xml:"FulfillmentDate,omitempty" json:"FulfillmentDate,omitempty"`
	ConsumerFulfillmentDate           *FulfillmentDate `xml:"ConsumerFulfillmentDate,omitempty" json:"ConsumerFulfillmentDate,omitempty"`
	FileAvailabilityDescription       *Description     `xml:"FileAvailabilityDescription,omitempty" json:"FileAvailabilityDescription,omitempty"`
	File                              *File            `xml:"File,omitempty" json:"File,omitempty"`
	Fingerprint                       *Fingerprint     `xml:"Fingerprint,omitempty" json:"Fingerprint,omitempty"`
	LanguageAndScriptCode             string           `xml:"LanguageAndScriptCode,attr,omitempty" json:"LanguageAndScriptCode,omitempty"`
}

type PreviewDetails struct {
	PartType          *Description `xml:"PartType,omitempty" json:"PartType,omitempty"`
	TopLeftCorner     string       `xml:"TopLeftCorner,omitempty" json:"TopLeftCorner,omitempty"`
	BottomRightCorner string       `xml:"BottomRightCorner,omitempty" json:"BottomRightCorner,omitempty"`
	ExpressionType    string       `xml:"ExpressionType,omitempty" json:"ExpressionType,omitempty"`
}

type ImageCodecType struct {
	Value            string `xml:",chardata"`
	Version          string `xml:"Version,attr,omitempty" json:"Version,omitempty"`
	Namespace        string `xml:"Namespace,attr,omitempty" json:"Namespace,omitempty"`
	UserDefinedValue string `xml:"UserDefinedValue,attr,omitempty" json:"UserDefinedValue,omitempty"`
}

type ImageType struct {
	Value            string `xml:",chardata"`
	Namespace        string `xml:"Namespace,attr,omitempty" json:"Namespace,omitempty"`
	UserDefinedValue string `xml:"UserDefinedValue,attr,omitempty" json:"UserDefinedValue,omitempty"`
}

type Video struct {
	VideoType                              *VideoType                              `xml:"VideoType,omitempty" json:"VideoType,omitempty"`
	IsArtistRelated                        bool                                    `xml:"IsArtistRelated,omitempty" json:"IsArtistRelated,omitempty"`
	VideoId                                *VideoId                                `xml:"VideoId,omitempty" json:"VideoId,omitempty"`
	IndirectVideoId                        *MusicalWorkId                          `xml:"IndirectVideoId,omitempty" json:"IndirectVideoId,omitempty"`
	ResourceReference                      string                                  `xml:"ResourceReference,omitempty" json:"ResourceReference,omitempty"`
	VideoCueSheetReference                 *VideoCueSheetReference                 `xml:"VideoCueSheetReference,omitempty" json:"VideoCueSheetReference,omitempty"`
	ReasonForCueSheetAbsence               *Reason                                 `xml:"ReasonForCueSheetAbsence,omitempty" json:"ReasonForCueSheetAbsence,omitempty"`
	ReferenceTitle                         *ReferenceTitle                         `xml:"ReferenceTitle,omitempty" json:"ReferenceTitle,omitempty"`
	Title                                  []*Title                                `xml:"Title,omitempty" json:"Title,omitempty"`
	InstrumentationDescription             *Description                            `xml:"InstrumentationDescription,omitempty" json:"InstrumentationDescription,omitempty"`
	IsMedley                               bool                                    `xml:"IsMedley,omitempty" json:"IsMedley,omitempty"`
	IsPotpourri                            bool                                    `xml:"IsPotpourri,omitempty" json:"IsPotpourri,omitempty"`
	IsInstrumental                         bool                                    `xml:"IsInstrumental,omitempty" json:"IsInstrumental,omitempty"`
	IsBackground                           bool                                    `xml:"IsBackground,omitempty" json:"IsBackground,omitempty"`
	IsHiddenResource                       bool                                    `xml:"IsHiddenResource,omitempty" json:"IsHiddenResource,omitempty"`
	IsBonusResource                        bool                                    `xml:"IsBonusResource,omitempty" json:"IsBonusResource,omitempty"`
	HasPreOrderFulfillment                 bool                                    `xml:"HasPreOrderFulfillment,omitempty" json:"HasPreOrderFulfillment,omitempty"`
	IsRemastered                           bool                                    `xml:"IsRemastered,omitempty" json:"IsRemastered,omitempty"`
	NoSilenceBefore                        bool                                    `xml:"NoSilenceBefore,omitempty" json:"NoSilenceBefore,omitempty"`
	NoSilenceAfter                         bool                                    `xml:"NoSilenceAfter,omitempty" json:"NoSilenceAfter,omitempty"`
	PerformerInformationRequired           bool                                    `xml:"PerformerInformationRequired,omitempty" json:"PerformerInformationRequired,omitempty"`
	LanguageOfPerformance                  string                                  `xml:"LanguageOfPerformance,omitempty" json:"LanguageOfPerformance,omitempty"`
	LanguageOfDubbing                      string                                  `xml:"LanguageOfDubbing,omitempty" json:"LanguageOfDubbing,omitempty"`
	SubTitleLanguage                       string                                  `xml:"SubTitleLanguage,omitempty" json:"SubTitleLanguage,omitempty"`
	Duration                               string                                  `xml:"Duration,omitempty" json:"Duration,omitempty"`
	RightsAgreementId                      *RightsAgreementId                      `xml:"RightsAgreementId,omitempty" json:"RightsAgreementId,omitempty"`
	VideoCollectionReferenceList           *SoundRecordingCollectionReferenceList  `xml:"VideoCollectionReferenceList,omitempty" json:"VideoCollectionReferenceList,omitempty"`
	ResourceMusicalWorkReferenceList       *ResourceMusicalWorkReferenceList       `xml:"ResourceMusicalWorkReferenceList,omitempty" json:"ResourceMusicalWorkReferenceList,omitempty"`
	ResourceContainedResourceReferenceList *ResourceContainedResourceReferenceList `xml:"ResourceContainedResourceReferenceList,omitempty" json:"ResourceContainedResourceReferenceList,omitempty"`
	CreationDate                           *EventDate                              `xml:"CreationDate,omitempty" json:"CreationDate,omitempty"`
	MasteredDate                           *EventDate                              `xml:"MasteredDate,omitempty" json:"MasteredDate,omitempty"`
	RemasteredDate                         *EventDate                              `xml:"RemasteredDate,omitempty" json:"RemasteredDate,omitempty"`
	VideoDetailsByTerritory                []*VideoDetailsByTerritory              `xml:"VideoDetailsByTerritory,omitempty" json:"VideoDetailsByTerritory,omitempty"`
	TerritoryOfCommissioning               *TerritoryCode                          `xml:"TerritoryOfCommissioning,omitempty" json:"TerritoryOfCommissioning,omitempty"`
	NumberOfFeaturedArtists                int                                     `xml:"NumberOfFeaturedArtists,omitempty" json:"NumberOfFeaturedArtists,omitempty"`
	NumberOfNonFeaturedArtists             int                                     `xml:"NumberOfNonFeaturedArtists,omitempty" json:"NumberOfNonFeaturedArtists,omitempty"`
	NumberOfContractedArtists              int                                     `xml:"NumberOfContractedArtists,omitempty" json:"NumberOfContractedArtists,omitempty"`
	NumberOfNonContractedArtists           int                                     `xml:"NumberOfNonContractedArtists,omitempty" json:"NumberOfNonContractedArtists,omitempty"`
	IsUpdated                              bool                                    `xml:"IsUpdated,attr,omitempty" json:"IsUpdated,omitempty"`
	LanguageAndScriptCode                  string                                  `xml:"LanguageAndScriptCode,attr,omitempty" json:"LanguageAndScriptCode,omitempty"`
}

type VideoDetailsByTerritory struct {
	TerritoryCode               *TerritoryCode                 `xml:"TerritoryCode,omitempty" json:"TerritoryCode,omitempty"`
	ExcludedTerritoryCode       *TerritoryCode                 `xml:"ExcludedTerritoryCode,omitempty" json:"ExcludedTerritoryCode,omitempty"`
	Title                       []*Title                       `xml:"Title,omitempty" json:"Title,omitempty"`
	DisplayArtist               []*Artist                      `xml:"DisplayArtist,omitempty" json:"DisplayArtist,omitempty"`
	DisplayConductor            []*Artist                      `xml:"DisplayConductor,omitempty" json:"DisplayConductor,omitempty"`
	ResourceContributor         []*DetailedResourceContributor `xml:"ResourceContributor,omitempty" json:"ResourceContributor,omitempty"`
	IndirectResourceContributor []*IndirectResourceContributor `xml:"IndirectResourceContributor,omitempty" json:"IndirectResourceContributor,omitempty"`
	RightsAgreementId           *RightsAgreementId             `xml:"RightsAgreementId,omitempty" json:"RightsAgreementId,omitempty"`
	DisplayArtistName           *Name                          `xml:"DisplayArtistName,omitempty" json:"DisplayArtistName,omitempty"`
	LabelName                   []*LabelName                   `xml:"LabelName,omitempty" json:"LabelName,omitempty"`
	RightsController            *TypedRightsController         `xml:"RightsController,omitempty" json:"RightsController,omitempty"`
	RemasteredDate              *EventDate                     `xml:"RemasteredDate,omitempty" json:"RemasteredDate,omitempty"`
	ResourceReleaseDate         *EventDate                     `xml:"ResourceReleaseDate,omitempty" json:"ResourceReleaseDate,omitempty"`
	OriginalResourceReleaseDate *EventDate                     `xml:"OriginalResourceReleaseDate,omitempty" json:"OriginalResourceReleaseDate,omitempty"`
	PLine                       *PLine                         `xml:"PLine,omitempty" json:"PLine,omitempty"`
	CourtesyLine                *CourtesyLine                  `xml:"CourtesyLine,omitempty" json:"CourtesyLine,omitempty"`
	SequenceNumber              int                            `xml:"SequenceNumber,omitempty" json:"SequenceNumber,omitempty"`
	HostSoundCarrier            *HostSoundCarrier              `xml:"HostSoundCarrier,omitempty" json:"HostSoundCarrier,omitempty"`
	MarketingComment            *Comment                       `xml:"MarketingComment,omitempty" json:"MarketingComment,omitempty"`
	Genre                       *Genre                         `xml:"Genre,omitempty" json:"Genre,omitempty"`
	ParentalWarningType         *ParentalWarningType           `xml:"ParentalWarningType,omitempty" json:"ParentalWarningType,omitempty"`
	AvRating                    *AvRating                      `xml:"AvRating,omitempty" json:"AvRating,omitempty"`
	FulfillmentDate             *FulfillmentDate               `xml:"FulfillmentDate,omitempty" json:"FulfillmentDate,omitempty"`
	Keywords                    *Keywords                      `xml:"Keywords,omitempty" json:"Keywords,omitempty"`
	Synopsis                    *Synopsis                      `xml:"Synopsis,omitempty" json:"Synopsis,omitempty"`
	CLine                       *CLine                         `xml:"CLine,omitempty" json:"CLine,omitempty"`
	TechnicalVideoDetails       *TechnicalVideoDetails         `xml:"TechnicalVideoDetails,omitempty" json:"TechnicalVideoDetails,omitempty"`
	Character                   *Character                     `xml:"Character,omitempty" json:"Character,omitempty"`
	LanguageAndScriptCode       string                         `xml:"LanguageAndScriptCode,attr,omitempty" json:"LanguageAndScriptCode,omitempty"`
}

type TechnicalVideoDetails struct {
	TechnicalResourceDetailsReference string                        `xml:"TechnicalResourceDetailsReference,omitempty" json:"TechnicalResourceDetailsReference,omitempty"`
	DrmPlatformType                   *DrmPlatformType              `xml:"DrmPlatformType,omitempty" json:"DrmPlatformType,omitempty"`
	OverallBitRate                    *BitRate                      `xml:"OverallBitRate,omitempty" json:"OverallBitRate,omitempty"`
	ContainerFormat                   *ContainerFormat              `xml:"ContainerFormat,omitempty" json:"ContainerFormat,omitempty"`
	VideoCodecType                    *VideoCodecType               `xml:"VideoCodecType,omitempty" json:"VideoCodecType,omitempty"`
	VideoBitRate                      *BitRate                      `xml:"VideoBitRate,omitempty" json:"VideoBitRate,omitempty"`
	FrameRate                         *FrameRate                    `xml:"FrameRate,omitempty" json:"FrameRate,omitempty"`
	ImageHeight                       *Extent                       `xml:"ImageHeight,omitempty" json:"ImageHeight,omitempty"`
	ImageWidth                        *Extent                       `xml:"ImageWidth,omitempty" json:"ImageWidth,omitempty"`
	AspectRatio                       *AspectRatio                  `xml:"AspectRatio,omitempty" json:"AspectRatio,omitempty"`
	ColorDepth                        int                           `xml:"ColorDepth,omitempty" json:"ColorDepth,omitempty"`
	VideoDefinitionType               string                        `xml:"VideoDefinitionType,omitempty" json:"VideoDefinitionType,omitempty"`
	AudioCodecType                    *AudioCodecType               `xml:"AudioCodecType,omitempty" json:"AudioCodecType,omitempty"`
	AudioBitRate                      *BitRate                      `xml:"AudioBitRate,omitempty" json:"AudioBitRate,omitempty"`
	NumberOfAudioChannels             int                           `xml:"NumberOfAudioChannels,omitempty" json:"NumberOfAudioChannels,omitempty"`
	AudioSamplingRate                 *SamplingRate                 `xml:"AudioSamplingRate,omitempty" json:"AudioSamplingRate,omitempty"`
	AudioBitsPerSample                int                           `xml:"AudioBitsPerSample,omitempty" json:"AudioBitsPerSample,omitempty"`
	Duration                          string                        `xml:"Duration,omitempty" json:"Duration,omitempty"`
	ResourceProcessingRequired        bool                          `xml:"ResourceProcessingRequired,omitempty" json:"ResourceProcessingRequired,omitempty"`
	UsableResourceDuration            string                        `xml:"UsableResourceDuration,omitempty" json:"UsableResourceDuration,omitempty"`
	IsPreview                         bool                          `xml:"IsPreview,omitempty" json:"IsPreview,omitempty"`
	PreviewDetails                    *SoundRecordingPreviewDetails `xml:"PreviewDetails,omitempty" json:"PreviewDetails,omitempty"`
	FulfillmentDate                   *FulfillmentDate              `xml:"FulfillmentDate,omitempty" json:"FulfillmentDate,omitempty"`
	ConsumerFulfillmentDate           *FulfillmentDate              `xml:"ConsumerFulfillmentDate,omitempty" json:"ConsumerFulfillmentDate,omitempty"`
	FileAvailabilityDescription       *Description                  `xml:"FileAvailabilityDescription,omitempty" json:"FileAvailabilityDescription,omitempty"`
	File                              *File                         `xml:"File,omitempty" json:"File,omitempty"`
	Fingerprint                       *Fingerprint                  `xml:"Fingerprint,omitempty" json:"Fingerprint,omitempty"`
	LanguageAndScriptCode             string                        `xml:"LanguageAndScriptCode,attr,omitempty" json:"LanguageAndScriptCode,omitempty"`
}

type VideoCodecType struct {
	Value            string `xml:",chardata"`
	Version          string `xml:"Version,attr,omitempty" json:"Version,omitempty"`
	Namespace        string `xml:"Namespace,attr,omitempty" json:"Namespace,omitempty"`
	UserDefinedValue string `xml:"UserDefinedValue,attr,omitempty" json:"UserDefinedValue,omitempty"`
}

type VideoType struct {
	Value            string `xml:",chardata"`
	Namespace        string `xml:"Namespace,attr,omitempty" json:"Namespace,omitempty"`
	UserDefinedValue string `xml:"UserDefinedValue,attr,omitempty" json:"UserDefinedValue,omitempty"`
}

type VideoId struct {
	ISRC          string         `xml:"ISRC,omitempty" json:"ISRC,omitempty"`
	ISAN          string         `xml:"ISAN,omitempty" json:"ISAN,omitempty"`
	VISAN         string         `xml:"VISAN,omitempty" json:"VISAN,omitempty"`
	CatalogNumber *CatalogNumber `xml:"CatalogNumber,omitempty" json:"CatalogNumber,omitempty"`
	ProprietaryId *ProprietaryId `xml:"ProprietaryId,omitempty" json:"ProprietaryId,omitempty"`
	EIDR          string         `xml:"EIDR,omitempty" json:"EIDR,omitempty"`
	IsReplaced    bool           `xml:"IsReplaced,attr,omitempty" json:"IsReplaced,omitempty"`
}

type VideoCueSheetReference struct {
	VideoCueSheetReference string `xml:"VideoCueSheetReference,omitempty" json:"VideoCueSheetReference,omitempty"`
}

type Reason struct {
	Value                 string `xml:",chardata"`
	LanguageAndScriptCode string `xml:"LanguageAndScriptCode,attr,omitempty" json:"LanguageAndScriptCode,omitempty"`
}

type MIDI struct {
	MidiType                               *MidiType                               `xml:"MidiType,omitempty" json:"MidiType,omitempty"`
	IsArtistRelated                        bool                                    `xml:"IsArtistRelated,omitempty" json:"IsArtistRelated,omitempty"`
	MidiId                                 *ResourceProprietaryId                  `xml:"MidiId,omitempty" json:"MidiId,omitempty"`
	IndirectMidiId                         *MusicalWorkId                          `xml:"IndirectMidiId,omitempty" json:"IndirectMidiId,omitempty"`
	ResourceReference                      string                                  `xml:"ResourceReference,omitempty" json:"ResourceReference,omitempty"`
	ReferenceTitle                         *ReferenceTitle                         `xml:"ReferenceTitle,omitempty" json:"ReferenceTitle,omitempty"`
	InstrumentationDescription             *Description                            `xml:"InstrumentationDescription,omitempty" json:"InstrumentationDescription,omitempty"`
	IsMedley                               bool                                    `xml:"IsMedley,omitempty" json:"IsMedley,omitempty"`
	IsPotpourri                            bool                                    `xml:"IsPotpourri,omitempty" json:"IsPotpourri,omitempty"`
	IsInstrumental                         bool                                    `xml:"IsInstrumental,omitempty" json:"IsInstrumental,omitempty"`
	IsBackground                           bool                                    `xml:"IsBackground,omitempty" json:"IsBackground,omitempty"`
	IsHiddenResource                       bool                                    `xml:"IsHiddenResource,omitempty" json:"IsHiddenResource,omitempty"`
	IsBonusResource                        bool                                    `xml:"IsBonusResource,omitempty" json:"IsBonusResource,omitempty"`
	IsComputerGenerated                    bool                                    `xml:"IsComputerGenerated,omitempty" json:"IsComputerGenerated,omitempty"`
	NoSilenceBefore                        bool                                    `xml:"NoSilenceBefore,omitempty" json:"NoSilenceBefore,omitempty"`
	NoSilenceAfter                         bool                                    `xml:"NoSilenceAfter,omitempty" json:"NoSilenceAfter,omitempty"`
	PerformerInformationRequired           bool                                    `xml:"PerformerInformationRequired,omitempty" json:"PerformerInformationRequired,omitempty"`
	LanguageOfPerformance                  string                                  `xml:"LanguageOfPerformance,omitempty" json:"LanguageOfPerformance,omitempty"`
	Duration                               string                                  `xml:"Duration,omitempty" json:"Duration,omitempty"`
	RightsAgreementId                      *RightsAgreementId                      `xml:"RightsAgreementId,omitempty" json:"RightsAgreementId,omitempty"`
	ResourceMusicalWorkReferenceList       *ResourceMusicalWorkReferenceList       `xml:"ResourceMusicalWorkReferenceList,omitempty" json:"ResourceMusicalWorkReferenceList,omitempty"`
	ResourceContainedResourceReferenceList *ResourceContainedResourceReferenceList `xml:"ResourceContainedResourceReferenceList,omitempty" json:"ResourceContainedResourceReferenceList,omitempty"`
	CreationDate                           *EventDate                              `xml:"CreationDate,omitempty" json:"CreationDate,omitempty"`
	MasteredDate                           *EventDate                              `xml:"MasteredDate,omitempty" json:"MasteredDate,omitempty"`
	RemasteredDate                         *EventDate                              `xml:"RemasteredDate,omitempty" json:"RemasteredDate,omitempty"`
	MidiDetailsByTerritory                 []*MidiDetailsByTerritory               `xml:"MidiDetailsByTerritory,omitempty" json:"MidiDetailsByTerritory,omitempty"`
	IsUpdated                              bool                                    `xml:"IsUpdated,attr,omitempty" json:"IsUpdated,omitempty"`
	LanguageAndScriptCode                  string                                  `xml:"LanguageAndScriptCode,attr,omitempty" json:"LanguageAndScriptCode,omitempty"`
}

type MidiDetailsByTerritory struct {
	TerritoryCode               *TerritoryCode                 `xml:"TerritoryCode,omitempty" json:"TerritoryCode,omitempty"`
	ExcludedTerritoryCode       *TerritoryCode                 `xml:"ExcludedTerritoryCode,omitempty" json:"ExcludedTerritoryCode,omitempty"`
	Title                       []*Title                       `xml:"Title,omitempty" json:"Title,omitempty"`
	DisplayArtist               []*Artist                      `xml:"DisplayArtist,omitempty" json:"DisplayArtist,omitempty"`
	ResourceContributor         []*DetailedResourceContributor `xml:"ResourceContributor,omitempty" json:"ResourceContributor,omitempty"`
	IndirectResourceContributor []*IndirectResourceContributor `xml:"IndirectResourceContributor,omitempty" json:"IndirectResourceContributor,omitempty"`
	RightsAgreementId           *RightsAgreementId             `xml:"RightsAgreementId,omitempty" json:"RightsAgreementId,omitempty"`
	DisplayArtistName           *Name                          `xml:"DisplayArtistName,omitempty" json:"DisplayArtistName,omitempty"`
	LabelName                   []*LabelName                   `xml:"LabelName,omitempty" json:"LabelName,omitempty"`
	RightsController            *TypedRightsController         `xml:"RightsController,omitempty" json:"RightsController,omitempty"`
	RemasteredDate              *EventDate                     `xml:"RemasteredDate,omitempty" json:"RemasteredDate,omitempty"`
	ResourceReleaseDate         *EventDate                     `xml:"ResourceReleaseDate,omitempty" json:"ResourceReleaseDate,omitempty"`
	OriginalResourceReleaseDate *EventDate                     `xml:"OriginalResourceReleaseDate,omitempty" json:"OriginalResourceReleaseDate,omitempty"`
	CLine                       *CLine                         `xml:"CLine,omitempty" json:"CLine,omitempty"`
	CourtesyLine                *CourtesyLine                  `xml:"CourtesyLine,omitempty" json:"CourtesyLine,omitempty"`
	SequenceNumber              int                            `xml:"SequenceNumber,omitempty" json:"SequenceNumber,omitempty"`
	HostSoundCarrier            *HostSoundCarrier              `xml:"HostSoundCarrier,omitempty" json:"HostSoundCarrier,omitempty"`
	MarketingComment            *Comment                       `xml:"MarketingComment,omitempty" json:"MarketingComment,omitempty"`
	Genre                       *Genre                         `xml:"Genre,omitempty" json:"Genre,omitempty"`
	ParentalWarningType         *ParentalWarningType           `xml:"ParentalWarningType,omitempty" json:"ParentalWarningType,omitempty"`
	FulfillmentDate             *FulfillmentDate               `xml:"FulfillmentDate,omitempty" json:"FulfillmentDate,omitempty"`
	Keywords                    *Keywords                      `xml:"Keywords,omitempty" json:"Keywords,omitempty"`
	Synopsis                    *Synopsis                      `xml:"Synopsis,omitempty" json:"Synopsis,omitempty"`
	TechnicalMidiDetails        *TechnicalMidiDetails          `xml:"TechnicalMidiDetails,omitempty" json:"TechnicalMidiDetails,omitempty"`
	LanguageAndScriptCode       string                         `xml:"LanguageAndScriptCode,attr,omitempty" json:"LanguageAndScriptCode,omitempty"`
}

type TechnicalMidiDetails struct {
	TechnicalResourceDetailsReference string                        `xml:"TechnicalResourceDetailsReference,omitempty" json:"TechnicalResourceDetailsReference,omitempty"`
	Duration                          string                        `xml:"Duration,omitempty" json:"Duration,omitempty"`
	ResourceProcessingRequired        bool                          `xml:"ResourceProcessingRequired,omitempty" json:"ResourceProcessingRequired,omitempty"`
	UsableResourceDuration            string                        `xml:"UsableResourceDuration,omitempty" json:"UsableResourceDuration,omitempty"`
	IsPreview                         bool                          `xml:"IsPreview,omitempty" json:"IsPreview,omitempty"`
	PreviewDetails                    *SoundRecordingPreviewDetails `xml:"PreviewDetails,omitempty" json:"PreviewDetails,omitempty"`
	FulfillmentDate                   *FulfillmentDate              `xml:"FulfillmentDate,omitempty" json:"FulfillmentDate,omitempty"`
	ConsumerFulfillmentDate           *FulfillmentDate              `xml:"ConsumerFulfillmentDate,omitempty" json:"ConsumerFulfillmentDate,omitempty"`
	FileAvailabilityDescription       *Description                  `xml:"FileAvailabilityDescription,omitempty" json:"FileAvailabilityDescription,omitempty"`
	File                              *File                         `xml:"File,omitempty" json:"File,omitempty"`
	NumberOfVoices                    int                           `xml:"NumberOfVoices,omitempty" json:"NumberOfVoices,omitempty"`
	SoundProcessorType                *SoundProcessorType           `xml:"SoundProcessorType,omitempty" json:"SoundProcessorType,omitempty"`
	Fingerprint                       *Fingerprint                  `xml:"Fingerprint,omitempty" json:"Fingerprint,omitempty"`
	LanguageAndScriptCode             string                        `xml:"LanguageAndScriptCode,attr,omitempty" json:"LanguageAndScriptCode,omitempty"`
}

type SoundProcessorType struct {
	Value            string `xml:",chardata"`
	Version          string `xml:"Version,attr,omitempty" json:"Version,omitempty"`
	Namespace        string `xml:"Namespace,attr,omitempty" json:"Namespace,omitempty"`
	UserDefinedValue string `xml:"UserDefinedValue,attr,omitempty" json:"UserDefinedValue,omitempty"`
}

type ResourceProprietaryId struct {
	ProprietaryId *ProprietaryId `xml:"ProprietaryId,omitempty" json:"ProprietaryId,omitempty"`
	IsReplaced    bool           `xml:"IsReplaced,attr,omitempty" json:"IsReplaced,omitempty"`
}

type MidiType struct {
	Value            string `xml:",chardata"`
	Namespace        string `xml:"Namespace,attr,omitempty" json:"Namespace,omitempty"`
	UserDefinedValue string `xml:"UserDefinedValue,attr,omitempty" json:"UserDefinedValue,omitempty"`
}

type SoundRecording struct {
	SoundRecordingType                     *SoundRecordingType                     `xml:"SoundRecordingType,omitempty" json:"SoundRecordingType,omitempty"`
	IsArtistRelated                        bool                                    `xml:"IsArtistRelated,omitempty" json:"IsArtistRelated,omitempty"`
	SoundRecordingId                       *SoundRecordingId                       `xml:"SoundRecordingId,omitempty" json:"SoundRecordingId,omitempty"`
	IndirectSoundRecordingId               *MusicalWorkId                          `xml:"IndirectSoundRecordingId,omitempty" json:"IndirectSoundRecordingId,omitempty"`
	ResourceReference                      string                                  `xml:"ResourceReference,omitempty" json:"ResourceReference,omitempty"`
	ReferenceTitle                         *ReferenceTitle                         `xml:"ReferenceTitle,omitempty" json:"ReferenceTitle,omitempty"`
	InstrumentationDescription             *Description                            `xml:"InstrumentationDescription,omitempty" json:"InstrumentationDescription,omitempty"`
	IsMedley                               bool                                    `xml:"IsMedley,omitempty" json:"IsMedley,omitempty"`
	IsPotpourri                            bool                                    `xml:"IsPotpourri,omitempty" json:"IsPotpourri,omitempty"`
	IsInstrumental                         bool                                    `xml:"IsInstrumental,omitempty" json:"IsInstrumental,omitempty"`
	IsBackground                           bool                                    `xml:"IsBackground,omitempty" json:"IsBackground,omitempty"`
	IsHiddenResource                       bool                                    `xml:"IsHiddenResource,omitempty" json:"IsHiddenResource,omitempty"`
	HasPreOrderFulfillment                 bool                                    `xml:"HasPreOrderFulfillment,omitempty" json:"HasPreOrderFulfillment,omitempty"`
	IsComputerGenerated                    bool                                    `xml:"IsComputerGenerated,omitempty" json:"IsComputerGenerated,omitempty"`
	IsRemastered                           bool                                    `xml:"IsRemastered,omitempty" json:"IsRemastered,omitempty"`
	NoSilenceBefore                        bool                                    `xml:"NoSilenceBefore,omitempty" json:"NoSilenceBefore,omitempty"`
	NoSilenceAfter                         bool                                    `xml:"NoSilenceAfter,omitempty" json:"NoSilenceAfter,omitempty"`
	PerformerInformationRequired           bool                                    `xml:"PerformerInformationRequired,omitempty" json:"PerformerInformationRequired,omitempty"`
	LanguageOfPerformance                  string                                  `xml:"LanguageOfPerformance,omitempty" json:"LanguageOfPerformance,omitempty"`
	Duration                               string                                  `xml:"Duration,omitempty" json:"Duration,omitempty"`
	RightsAgreementId                      *RightsAgreementId                      `xml:"RightsAgreementId,omitempty" json:"RightsAgreementId,omitempty"`
	SoundRecordingCollectionReferenceList  *SoundRecordingCollectionReferenceList  `xml:"SoundRecordingCollectionReferenceList,omitempty" json:"SoundRecordingCollectionReferenceList,omitempty"`
	ResourceMusicalWorkReferenceList       *ResourceMusicalWorkReferenceList       `xml:"ResourceMusicalWorkReferenceList,omitempty" json:"ResourceMusicalWorkReferenceList,omitempty"`
	ResourceContainedResourceReferenceList *ResourceContainedResourceReferenceList `xml:"ResourceContainedResourceReferenceList,omitempty" json:"ResourceContainedResourceReferenceList,omitempty"`
	CreationDate                           *EventDate                              `xml:"CreationDate,omitempty" json:"CreationDate,omitempty"`
	MasteredDate                           *EventDate                              `xml:"MasteredDate,omitempty" json:"MasteredDate,omitempty"`
	RemasteredDate                         *EventDate                              `xml:"RemasteredDate,omitempty" json:"RemasteredDate,omitempty"`
	SoundRecordingDetailsByTerritory       []*SoundRecordingDetailsByTerritory     `xml:"SoundRecordingDetailsByTerritory,omitempty" json:"SoundRecordingDetailsByTerritory,omitempty"`
	TerritoryOfCommissioning               *TerritoryCode                          `xml:"TerritoryOfCommissioning,omitempty" json:"TerritoryOfCommissioning,omitempty"`
	NumberOfFeaturedArtists                int                                     `xml:"NumberOfFeaturedArtists,omitempty" json:"NumberOfFeaturedArtists,omitempty"`
	NumberOfNonFeaturedArtists             int                                     `xml:"NumberOfNonFeaturedArtists,omitempty" json:"NumberOfNonFeaturedArtists,omitempty"`
	NumberOfContractedArtists              int                                     `xml:"NumberOfContractedArtists,omitempty" json:"NumberOfContractedArtists,omitempty"`
	NumberOfNonContractedArtists           int                                     `xml:"NumberOfNonContractedArtists,omitempty" json:"NumberOfNonContractedArtists,omitempty"`
	IsUpdated                              bool                                    `xml:"IsUpdated,attr,omitempty" json:"IsUpdated,omitempty"`
	LanguageAndScriptCode                  string                                  `xml:"LanguageAndScriptCode,attr,omitempty" json:"LanguageAndScriptCode,omitempty"`
}

type SoundRecordingDetailsByTerritory struct {
	TerritoryCode                  *TerritoryCode                  `xml:"TerritoryCode,omitempty" json:"TerritoryCode,omitempty"`
	ExcludedTerritoryCode          *TerritoryCode                  `xml:"ExcludedTerritoryCode,omitempty" json:"ExcludedTerritoryCode,omitempty"`
	Title                          []*Title                        `xml:"Title,omitempty" json:"Title,omitempty"`
	DisplayArtist                  []*Artist                       `xml:"DisplayArtist,omitempty" json:"DisplayArtist,omitempty"`
	DisplayConductor               []*Artist                       `xml:"DisplayConductor,omitempty" json:"DisplayConductor,omitempty"`
	ResourceContributor            []*DetailedResourceContributor  `xml:"ResourceContributor,omitempty" json:"ResourceContributor,omitempty"`
	IndirectResourceContributor    []*IndirectResourceContributor  `xml:"IndirectResourceContributor,omitempty" json:"IndirectResourceContributor,omitempty"`
	RightsAgreementId              *RightsAgreementId              `xml:"RightsAgreementId,omitempty" json:"RightsAgreementId,omitempty"`
	DisplayArtistName              *Name                           `xml:"DisplayArtistName,omitempty" json:"DisplayArtistName,omitempty"`
	LabelName                      []*LabelName                    `xml:"LabelName,omitempty" json:"LabelName,omitempty"`
	RightsController               *TypedRightsController          `xml:"RightsController,omitempty" json:"RightsController,omitempty"`
	RemasteredDate                 *EventDate                      `xml:"RemasteredDate,omitempty" json:"RemasteredDate,omitempty"`
	ResourceReleaseDate            *EventDate                      `xml:"ResourceReleaseDate,omitempty" json:"ResourceReleaseDate,omitempty"`
	OriginalResourceReleaseDate    *EventDate                      `xml:"OriginalResourceReleaseDate,omitempty" json:"OriginalResourceReleaseDate,omitempty"`
	PLine                          *PLine                          `xml:"PLine,omitempty" json:"PLine,omitempty"`
	CourtesyLine                   *CourtesyLine                   `xml:"CourtesyLine,omitempty" json:"CourtesyLine,omitempty"`
	SequenceNumber                 int                             `xml:"SequenceNumber,omitempty" json:"SequenceNumber,omitempty"`
	HostSoundCarrier               *HostSoundCarrier               `xml:"HostSoundCarrier,omitempty" json:"HostSoundCarrier,omitempty"`
	MarketingComment               *Comment                        `xml:"MarketingComment,omitempty" json:"MarketingComment,omitempty"`
	Genre                          *Genre                          `xml:"Genre,omitempty" json:"Genre,omitempty"`
	ParentalWarningType            *ParentalWarningType            `xml:"ParentalWarningType,omitempty" json:"ParentalWarningType,omitempty"`
	AvRating                       *AvRating                       `xml:"AvRating,omitempty" json:"AvRating,omitempty"`
	TechnicalSoundRecordingDetails *TechnicalSoundRecordingDetails `xml:"TechnicalSoundRecordingDetails,omitempty" json:"TechnicalSoundRecordingDetails,omitempty"`
	FulfillmentDate                *FulfillmentDate                `xml:"FulfillmentDate,omitempty" json:"FulfillmentDate,omitempty"`
	Keywords                       *Keywords                       `xml:"Keywords,omitempty" json:"Keywords,omitempty"`
	Synopsis                       *Synopsis                       `xml:"Synopsis,omitempty" json:"Synopsis,omitempty"`
	LanguageAndScriptCode          string                          `xml:"LanguageAndScriptCode,attr,omitempty" json:"LanguageAndScriptCode,omitempty"`
}

type Keywords struct {
	Value                 string `xml:",chardata"`
	LanguageAndScriptCode string `xml:"LanguageAndScriptCode,attr,omitempty" json:"LanguageAndScriptCode,omitempty"`
}

type Synopsis struct {
	Value                 string `xml:",chardata"`
	LanguageAndScriptCode string `xml:"LanguageAndScriptCode,attr,omitempty" json:"LanguageAndScriptCode,omitempty"`
}

type TechnicalSoundRecordingDetails struct {
	TechnicalResourceDetailsReference string                        `xml:"TechnicalResourceDetailsReference,omitempty" json:"TechnicalResourceDetailsReference,omitempty"`
	DrmPlatformType                   *DrmPlatformType              `xml:"DrmPlatformType,omitempty" json:"DrmPlatformType,omitempty"`
	ContainerFormat                   *ContainerFormat              `xml:"ContainerFormat,omitempty" json:"ContainerFormat,omitempty"`
	AudioCodecType                    *AudioCodecType               `xml:"AudioCodecType,omitempty" json:"AudioCodecType,omitempty"`
	BitRate                           *BitRate                      `xml:"BitRate,omitempty" json:"BitRate,omitempty"`
	NumberOfChannels                  int                           `xml:"NumberOfChannels,omitempty" json:"NumberOfChannels,omitempty"`
	SamplingRate                      *SamplingRate                 `xml:"SamplingRate,omitempty" json:"SamplingRate,omitempty"`
	BitsPerSample                     int                           `xml:"BitsPerSample,omitempty" json:"BitsPerSample,omitempty"`
	Duration                          string                        `xml:"Duration,omitempty" json:"Duration,omitempty"`
	ResourceProcessingRequired        bool                          `xml:"ResourceProcessingRequired,omitempty" json:"ResourceProcessingRequired,omitempty"`
	UsableResourceDuration            string                        `xml:"UsableResourceDuration,omitempty" json:"UsableResourceDuration,omitempty"`
	IsPreview                         bool                          `xml:"IsPreview,omitempty" json:"IsPreview,omitempty"`
	PreviewDetails                    *SoundRecordingPreviewDetails `xml:"PreviewDetails,omitempty" json:"PreviewDetails,omitempty"`
	FulfillmentDate                   *FulfillmentDate              `xml:"FulfillmentDate,omitempty" json:"FulfillmentDate,omitempty"`
	ConsumerFulfillmentDate           *FulfillmentDate              `xml:"ConsumerFulfillmentDate,omitempty" json:"ConsumerFulfillmentDate,omitempty"`
	FileAvailabilityDescription       *Description                  `xml:"FileAvailabilityDescription,omitempty" json:"FileAvailabilityDescription,omitempty"`
	File                              *File                         `xml:"File,omitempty" json:"File,omitempty"`
	Fingerprint                       *Fingerprint                  `xml:"Fingerprint,omitempty" json:"Fingerprint,omitempty"`
	LanguageAndScriptCode             string                        `xml:"LanguageAndScriptCode,attr,omitempty" json:"LanguageAndScriptCode,omitempty"`
}

type Fingerprint struct {
	Fingerprint                   string                    `xml:"Fingerprint,omitempty" json:"Fingerprint,omitempty"`
	FingerprintAlgorithmType      *FingerprintAlgorithmType `xml:"FingerprintAlgorithmType,omitempty" json:"FingerprintAlgorithmType,omitempty"`
	FingerprintAlgorithmVersion   string                    `xml:"FingerprintAlgorithmVersion,omitempty" json:"FingerprintAlgorithmVersion,omitempty"`
	FingerprintAlgorithmParameter string                    `xml:"FingerprintAlgorithmParameter,omitempty" json:"FingerprintAlgorithmParameter,omitempty"`
	FingerprintDataType           string                    `xml:"FingerprintDataType,omitempty" json:"FingerprintDataType,omitempty"`
}

type FingerprintAlgorithmType struct {
	Value            string `xml:",chardata"`
	Namespace        string `xml:"Namespace,attr,omitempty" json:"Namespace,omitempty"`
	UserDefinedValue string `xml:"UserDefinedValue,attr,omitempty" json:"UserDefinedValue,omitempty"`
}

type File struct {
	FileName string   `xml:"FileName,omitempty" json:"FileName,omitempty"`
	FilePath string   `xml:"FilePath,omitempty" json:"FilePath,omitempty"`
	URL      string   `xml:"URL,omitempty" json:"URL,omitempty"`
	HashSum  *HashSum `xml:"HashSum,omitempty" json:"HashSum,omitempty"`
}

type HashSum struct {
	HashSum              string                `xml:"HashSum,omitempty" json:"HashSum,omitempty"`
	HashSumAlgorithmType *HashSumAlgorithmType `xml:"HashSumAlgorithmType,omitempty" json:"HashSumAlgorithmType,omitempty"`
	HashSumDataType      string                `xml:"HashSumDataType,omitempty" json:"HashSumDataType,omitempty"`
}

type HashSumAlgorithmType struct {
	Value            string `xml:",chardata"`
	Namespace        string `xml:"Namespace,attr,omitempty" json:"Namespace,omitempty"`
	UserDefinedValue string `xml:"UserDefinedValue,attr,omitempty" json:"UserDefinedValue,omitempty"`
}

type FulfillmentDate struct {
	FulfillmentDate          string `xml:"FulfillmentDate,omitempty" json:"FulfillmentDate,omitempty"`
	ResourceReleaseReference string `xml:"ResourceReleaseReference,omitempty" json:"ResourceReleaseReference,omitempty"`
}

type SoundRecordingPreviewDetails struct {
	PartType          *Description `xml:"PartType,omitempty" json:"PartType,omitempty"`
	StartPoint        string       `xml:"StartPoint,omitempty" json:"StartPoint,omitempty"`
	EndPoint          string       `xml:"EndPoint,omitempty" json:"EndPoint,omitempty"`
	Duration          string       `xml:"Duration,omitempty" json:"Duration,omitempty"`
	TopLeftCorner     string       `xml:"TopLeftCorner,omitempty" json:"TopLeftCorner,omitempty"`
	BottomRightCorner string       `xml:"BottomRightCorner,omitempty" json:"BottomRightCorner,omitempty"`
	ExpressionType    string       `xml:"ExpressionType,omitempty" json:"ExpressionType,omitempty"`
}

type DrmPlatformType struct {
	Value            string `xml:",chardata"`
	Version          string `xml:"Version,attr,omitempty" json:"Version,omitempty"`
	Namespace        string `xml:"Namespace,attr,omitempty" json:"Namespace,omitempty"`
	UserDefinedValue string `xml:"UserDefinedValue,attr,omitempty" json:"UserDefinedValue,omitempty"`
}

type ContainerFormat struct {
	Value            string `xml:",chardata"`
	Namespace        string `xml:"Namespace,attr,omitempty" json:"Namespace,omitempty"`
	UserDefinedValue string `xml:"UserDefinedValue,attr,omitempty" json:"UserDefinedValue,omitempty"`
}

type AudioCodecType struct {
	Value            string `xml:",chardata"`
	Version          string `xml:"Version,attr,omitempty" json:"Version,omitempty"`
	Namespace        string `xml:"Namespace,attr,omitempty" json:"Namespace,omitempty"`
	UserDefinedValue string `xml:"UserDefinedValue,attr,omitempty" json:"UserDefinedValue,omitempty"`
}

type BitRate struct {
	Value         string `xml:",chardata"`
	UnitOfMeasure string `xml:"UnitOfMeasure,attr,omitempty" json:"UnitOfMeasure,omitempty"`
}

type FrameRate struct {
	Value         string `xml:",chardata"`
	UnitOfMeasure string `xml:"UnitOfMeasure,attr,omitempty" json:"UnitOfMeasure,omitempty"`
}

type Extent struct {
	Value         string `xml:",chardata"`
	UnitOfMeasure string `xml:"UnitOfMeasure,attr,omitempty" json:"UnitOfMeasure,omitempty"`
}

type AspectRatio struct {
	Value           string `xml:",chardata"`
	AspectRatioType string `xml:"AspectRatioType,attr,omitempty" json:"AspectRatioType,omitempty"`
}

type SamplingRate struct {
	Value         string `xml:",chardata"`
	UnitOfMeasure string `xml:"UnitOfMeasure,attr,omitempty" json:"UnitOfMeasure,omitempty"`
}

type AvRating struct {
	RatingText              string        `xml:"RatingText,omitempty" json:"RatingText,omitempty"`
	RatingAgency            *RatingAgency `xml:"RatingAgency,omitempty" json:"RatingAgency,omitempty"`
	RatingSchemeDescription *Description  `xml:"RatingSchemeDescription,omitempty" json:"RatingSchemeDescription,omitempty"`
}

type RatingAgency struct {
	Value            string `xml:",chardata"`
	Namespace        string `xml:"Namespace,attr,omitempty" json:"Namespace,omitempty"`
	UserDefinedValue string `xml:"UserDefinedValue,attr,omitempty" json:"UserDefinedValue,omitempty"`
}

type ParentalWarningType struct {
	Value            string `xml:",chardata"`
	Namespace        string `xml:"Namespace,attr,omitempty" json:"Namespace,omitempty"`
	UserDefinedValue string `xml:"UserDefinedValue,attr,omitempty" json:"UserDefinedValue,omitempty"`
}

type HostSoundCarrier struct {
	ReleaseId                   *ReleaseId                   `xml:"ReleaseId,omitempty" json:"ReleaseId,omitempty"`
	RightsAgreementId           *RightsAgreementId           `xml:"RightsAgreementId,omitempty" json:"RightsAgreementId,omitempty"`
	Title                       []*Title                     `xml:"Title,omitempty" json:"Title,omitempty"`
	DisplayArtist               []*Artist                    `xml:"DisplayArtist,omitempty" json:"DisplayArtist,omitempty"`
	AdministratingRecordCompany *AdministratingRecordCompany `xml:"AdministratingRecordCompany,omitempty" json:"AdministratingRecordCompany,omitempty"`
	TrackNumber                 string                       `xml:"TrackNumber,omitempty" json:"TrackNumber,omitempty"`
	VolumeNumberInSet           string                       `xml:"VolumeNumberInSet,omitempty" json:"VolumeNumberInSet,omitempty"`
}

type ReleaseId struct {
	GRid          string         `xml:"GRid,omitempty" json:"GRid,omitempty"`
	ISRC          string         `xml:"ISRC,omitempty" json:"ISRC,omitempty"`
	ICPN          *ICPN          `xml:"ICPN,omitempty" json:"ICPN,omitempty"`
	CatalogNumber *CatalogNumber `xml:"CatalogNumber,omitempty" json:"CatalogNumber,omitempty"`
	ProprietaryId *ProprietaryId `xml:"ProprietaryId,omitempty" json:"ProprietaryId,omitempty"`
	IsReplaced    bool           `xml:"IsReplaced,attr,omitempty" json:"IsReplaced,omitempty"`
}

type ICPN struct {
	Value string `xml:",chardata"`
	IsEan bool   `xml:"IsEan,attr,omitempty" json:"IsEan,omitempty"`
}

type AdministratingRecordCompany struct {
	PartyId          *PartyId   `xml:"PartyId,omitempty" json:"PartyId,omitempty"`
	PartyName        *PartyName `xml:"PartyName,omitempty" json:"PartyName,omitempty"`
	Namespace        string     `xml:"Namespace,attr,omitempty" json:"Namespace,omitempty"`
	UserDefinedValue string     `xml:"UserDefinedValue,attr,omitempty" json:"UserDefinedValue,omitempty"`
	Role             string     `xml:"Role,attr,omitempty" json:"Role,omitempty"`
}

type CourtesyLine struct {
	Value                 string `xml:",chardata"`
	LanguageAndScriptCode string `xml:"LanguageAndScriptCode,attr,omitempty" json:"LanguageAndScriptCode,omitempty"`
}

type LabelName struct {
	Value                 string `xml:",chardata"`
	LabelNameType         string `xml:"LabelNameType,omitempty" json:"LabelNameType,omitempty"`
	Namespace             string `xml:"Namespace,attr,omitempty" json:"Namespace,omitempty"`
	UserDefinedValue      string `xml:"UserDefinedValue,attr,omitempty" json:"UserDefinedValue,omitempty"`
	LanguageAndScriptCode string `xml:"LanguageAndScriptCode,attr,omitempty" json:"LanguageAndScriptCode,omitempty"`
}

type IndirectResourceContributor struct {
	PartyId                         *PartyId                      `xml:"PartyId,omitempty" json:"PartyId,omitempty"`
	PartyName                       *PartyName                    `xml:"PartyName,omitempty" json:"PartyName,omitempty"`
	IndirectResourceContributorRole []*MusicalWorkContributorRole `xml:"IndirectResourceContributorRole,omitempty" json:"IndirectResourceContributorRole,omitempty"`
	Nationality                     *TerritoryCode                `xml:"Nationality,omitempty" json:"Nationality,omitempty"`
	SequenceNumber                  int                           `xml:"SequenceNumber,attr,omitempty" json:"SequenceNumber,omitempty"`
}

type Artist struct {
	PartyId        *PartyId       `xml:"PartyId,omitempty" json:"PartyId,omitempty"`
	PartyName      *PartyName     `xml:"PartyName,omitempty" json:"PartyName,omitempty"`
	ArtistRole     *ArtistRole    `xml:"ArtistRole,omitempty" json:"ArtistRole,omitempty"`
	Nationality    *TerritoryCode `xml:"Nationality,omitempty" json:"Nationality,omitempty"`
	SequenceNumber int            `xml:"SequenceNumber,attr,omitempty" json:"SequenceNumber,omitempty"`
}

type ResourceContainedResourceReferenceList struct {
	ResourceContainedResourceReference []*ResourceContainedResourceReference `xml:"ResourceContainedResourceReference,omitempty" json:"ResourceContainedResourceReference,omitempty"`
}

type ResourceContainedResourceReference struct {
	ResourceContainedResourceReference string   `xml:"ResourceContainedResourceReference,omitempty" json:"ResourceContainedResourceReference,omitempty"`
	DurationUsed                       string   `xml:"DurationUsed,omitempty" json:"DurationUsed,omitempty"`
	StartPoint                         string   `xml:"StartPoint,omitempty" json:"StartPoint,omitempty"`
	Purpose                            *Purpose `xml:"Purpose,omitempty" json:"Purpose,omitempty"`
}

type Purpose struct {
	Value            string `xml:",chardata"`
	Namespace        string `xml:"Namespace,attr,omitempty" json:"Namespace,omitempty"`
	UserDefinedValue string `xml:"UserDefinedValue,attr,omitempty" json:"UserDefinedValue,omitempty"`
}

type ResourceMusicalWorkReferenceList struct {
	ResourceMusicalWorkReference []*ResourceMusicalWorkReference `xml:"ResourceMusicalWorkReference,omitempty" json:"ResourceMusicalWorkReference,omitempty"`
}

type ResourceMusicalWorkReference struct {
	SequenceNumber               int    `xml:"SequenceNumber,omitempty" json:"SequenceNumber,omitempty"`
	DurationUsed                 string `xml:"DurationUsed,omitempty" json:"DurationUsed,omitempty"`
	IsFragment                   bool   `xml:"IsFragment,omitempty" json:"IsFragment,omitempty"`
	ResourceMusicalWorkReference string `xml:"ResourceMusicalWorkReference,omitempty" json:"ResourceMusicalWorkReference,omitempty"`
}

type SoundRecordingCollectionReferenceList struct {
	NumberOfCollections               int                                  `xml:"NumberOfCollections,omitempty" json:"NumberOfCollections,omitempty"`
	SoundRecordingCollectionReference []*SoundRecordingCollectionReference `xml:"SoundRecordingCollectionReference,omitempty" json:"SoundRecordingCollectionReference,omitempty"`
}

type SoundRecordingCollectionReference struct {
	SequenceNumber                    int    `xml:"SequenceNumber,omitempty" json:"SequenceNumber,omitempty"`
	SoundRecordingCollectionReference string `xml:"SoundRecordingCollectionReference,omitempty" json:"SoundRecordingCollectionReference,omitempty"`
	StartTime                         string `xml:"StartTime,omitempty" json:"StartTime,omitempty"`
	Duration                          string `xml:"Duration,omitempty" json:"Duration,omitempty"`
	EndTime                           string `xml:"EndTime,omitempty" json:"EndTime,omitempty"`
	ReleaseResourceType               string `xml:"ReleaseResourceType,omitempty" json:"ReleaseResourceType,omitempty"`
}

type SoundRecordingType struct {
	Value            string `xml:",chardata"`
	Namespace        string `xml:"Namespace,attr,omitempty" json:"Namespace,omitempty"`
	UserDefinedValue string `xml:"UserDefinedValue,attr,omitempty" json:"UserDefinedValue,omitempty"`
}

type SoundRecordingId struct {
	ISRC          string         `xml:"ISRC,omitempty" json:"ISRC,omitempty"`
	CatalogNumber *CatalogNumber `xml:"CatalogNumber,omitempty" json:"CatalogNumber,omitempty"`
	ProprietaryId *ProprietaryId `xml:"ProprietaryId,omitempty" json:"ProprietaryId,omitempty"`
	IsReplaced    bool           `xml:"IsReplaced,attr,omitempty" json:"IsReplaced,omitempty"`
}

type CueSheetList struct {
	CueSheet []*CueSheet `xml:"CueSheet,omitempty" json:"CueSheet,omitempty"`
}

type CueSheet struct {
	CueSheetId        *ProprietaryId `xml:"CueSheetId,omitempty" json:"CueSheetId,omitempty"`
	CueSheetReference string         `xml:"CueSheetReference,omitempty" json:"CueSheetReference,omitempty"`
	CueSheetType      *CueSheetType  `xml:"CueSheetType,omitempty" json:"CueSheetType,omitempty"`
	Cue               *Cue           `xml:"Cue,omitempty" json:"Cue,omitempty"`
}

type Cue struct {
	CueUseType                            *CueUseType                    `xml:"CueUseType,omitempty" json:"CueUseType,omitempty"`
	CueThemeType                          *CueThemeType                  `xml:"CueThemeType,omitempty" json:"CueThemeType,omitempty"`
	CueVocalType                          *CueVocalType                  `xml:"CueVocalType,omitempty" json:"CueVocalType,omitempty"`
	IsDance                               bool                           `xml:"IsDance,omitempty" json:"IsDance,omitempty"`
	CueVisualPerceptionType               *CueVisualPerceptionType       `xml:"CueVisualPerceptionType,omitempty" json:"CueVisualPerceptionType,omitempty"`
	CueOrigin                             *CueOrigin                     `xml:"CueOrigin,omitempty" json:"CueOrigin,omitempty"`
	CueCreationReference                  *CueCreationReference          `xml:"CueCreationReference,omitempty" json:"CueCreationReference,omitempty"`
	ReferencedCreationType                string                         `xml:"ReferencedCreationType,omitempty" json:"ReferencedCreationType,omitempty"`
	ReferencedCreationId                  *CreationId                    `xml:"ReferencedCreationId,omitempty" json:"ReferencedCreationId,omitempty"`
	ReferencedCreationTitle               []*Title                       `xml:"ReferencedCreationTitle,omitempty" json:"ReferencedCreationTitle,omitempty"`
	ReferencedCreationContributor         []*DetailedResourceContributor `xml:"ReferencedCreationContributor,omitempty" json:"ReferencedCreationContributor,omitempty"`
	ReferencedIndirectCreationContributor []*MusicalWorkContributor      `xml:"ReferencedIndirectCreationContributor,omitempty" json:"ReferencedIndirectCreationContributor,omitempty"`
	ReferencedCreationCharacter           *Character                     `xml:"ReferencedCreationCharacter,omitempty" json:"ReferencedCreationCharacter,omitempty"`
	HasMusicalContent                     bool                           `xml:"HasMusicalContent,omitempty" json:"HasMusicalContent,omitempty"`
	StartTime                             string                         `xml:"StartTime,omitempty" json:"StartTime,omitempty"`
	Duration                              string                         `xml:"Duration,omitempty" json:"Duration,omitempty"`
	EndTime                               string                         `xml:"EndTime,omitempty" json:"EndTime,omitempty"`
	PLine                                 *PLine                         `xml:"PLine,omitempty" json:"PLine,omitempty"`
	CLine                                 *CLine                         `xml:"CLine,omitempty" json:"CLine,omitempty"`
}

type PLine struct {
	Year                  string `xml:"Year,omitempty" json:"Year,omitempty"`
	PLineCompany          string `xml:"PLineCompany,omitempty" json:"PLineCompany,omitempty"`
	PLineText             string `xml:"PLineText,omitempty" json:"PLineText,omitempty"`
	LanguageAndScriptCode string `xml:"LanguageAndScriptCode,attr,omitempty" json:"LanguageAndScriptCode,omitempty"`
}

type CLine struct {
	Year                  string `xml:"Year,omitempty" json:"Year,omitempty"`
	CLineCompany          string `xml:"CLineCompany,omitempty" json:"CLineCompany,omitempty"`
	CLineText             string `xml:"CLineText,omitempty" json:"CLineText,omitempty"`
	LanguageAndScriptCode string `xml:"LanguageAndScriptCode,attr,omitempty" json:"LanguageAndScriptCode,omitempty"`
}

type Character struct {
	PartyId             *PartyId                       `xml:"PartyId,omitempty" json:"PartyId,omitempty"`
	PartyName           *PartyName                     `xml:"PartyName,omitempty" json:"PartyName,omitempty"`
	ResourceContributor []*DetailedResourceContributor `xml:"ResourceContributor,omitempty" json:"ResourceContributor,omitempty"`
	SequenceNumber      int                            `xml:"SequenceNumber,attr,omitempty" json:"SequenceNumber,omitempty"`
}

type DetailedResourceContributor struct {
	PartyId                    *PartyId                    `xml:"PartyId,omitempty" json:"PartyId,omitempty"`
	PartyName                  *PartyName                  `xml:"PartyName,omitempty" json:"PartyName,omitempty"`
	ResourceContributorRole    []*ResourceContributorRole  `xml:"ResourceContributorRole,omitempty" json:"ResourceContributorRole,omitempty"`
	IsFeaturedArtist           bool                        `xml:"IsFeaturedArtist,omitempty" json:"IsFeaturedArtist,omitempty"`
	IsContractedArtist         bool                        `xml:"IsContractedArtist,omitempty" json:"IsContractedArtist,omitempty"`
	InstrumentType             string                      `xml:"InstrumentType,omitempty" json:"InstrumentType,omitempty"`
	ArtistDelegatedUsageRights *ArtistDelegatedUsageRights `xml:"ArtistDelegatedUsageRights,omitempty" json:"ArtistDelegatedUsageRights,omitempty"`
	Sex                        string                      `xml:"Sex,omitempty" json:"Sex,omitempty"`
	Nationality                *TerritoryCode              `xml:"Nationality,omitempty" json:"Nationality,omitempty"`
	DateAndPlaceOfBirth        *EventDate                  `xml:"DateAndPlaceOfBirth,omitempty" json:"DateAndPlaceOfBirth,omitempty"`
	DateAndPlaceOfDeath        *EventDate                  `xml:"DateAndPlaceOfDeath,omitempty" json:"DateAndPlaceOfDeath,omitempty"`
	PrimaryRole                *ArtistRole                 `xml:"PrimaryRole,omitempty" json:"PrimaryRole,omitempty"`
	Performance                *Performance                `xml:"Performance,omitempty" json:"Performance,omitempty"`
	PrimaryInstrumentType      string                      `xml:"PrimaryInstrumentType,omitempty" json:"PrimaryInstrumentType,omitempty"`
	GoverningAgreementType     *GoverningAgreementType     `xml:"GoverningAgreementType,omitempty" json:"GoverningAgreementType,omitempty"`
	ContactInformation         *ContactId                  `xml:"ContactInformation,omitempty" json:"ContactInformation,omitempty"`
	TerritoryOfResidency       *TerritoryCode              `xml:"TerritoryOfResidency,omitempty" json:"TerritoryOfResidency,omitempty"`
	Citizenship                *TerritoryCode              `xml:"Citizenship,omitempty" json:"Citizenship,omitempty"`
	AdditionalRoles            []*ArtistRole               `xml:"AdditionalRoles,omitempty" json:"AdditionalRoles,omitempty"`
	Genre                      *Genre                      `xml:"Genre,omitempty" json:"Genre,omitempty"`
	Membership                 *Membership                 `xml:"Membership,omitempty" json:"Membership,omitempty"`
	SequenceNumber             int                         `xml:"SequenceNumber,attr,omitempty" json:"SequenceNumber,omitempty"`
}

type ResourceContributorRole struct {
	Value            string `xml:",chardata"`
	Namespace        string `xml:"Namespace,attr,omitempty" json:"Namespace,omitempty"`
	UserDefinedValue string `xml:"UserDefinedValue,attr,omitempty" json:"UserDefinedValue,omitempty"`
}

type Membership struct {
	Organization   *PartyDescriptor `xml:"Organization,omitempty" json:"Organization,omitempty"`
	MembershipType string           `xml:"MembershipType,omitempty" json:"MembershipType,omitempty"`
	StartDate      string           `xml:"StartDate,omitempty" json:"StartDate,omitempty"`
	EndDate        string           `xml:"EndDate,omitempty" json:"EndDate,omitempty"`
}

type PartyDescriptor struct {
	PartyId   *PartyId   `xml:"PartyId,omitempty" json:"PartyId,omitempty"`
	PartyName *PartyName `xml:"PartyName,omitempty" json:"PartyName,omitempty"`
}

type Genre struct {
	GenreText             *Description `xml:"GenreText,omitempty" json:"GenreText,omitempty"`
	SubGenre              *Description `xml:"SubGenre,omitempty" json:"SubGenre,omitempty"`
	LanguageAndScriptCode string       `xml:"LanguageAndScriptCode,attr,omitempty" json:"LanguageAndScriptCode,omitempty"`
}

type Description struct {
	Value                 string `xml:",chardata"`
	LanguageAndScriptCode string `xml:"LanguageAndScriptCode,attr,omitempty" json:"LanguageAndScriptCode,omitempty"`
}

type Comment struct {
	Value                 string `xml:",chardata"`
	LanguageAndScriptCode string `xml:"LanguageAndScriptCode,attr,omitempty" json:"LanguageAndScriptCode,omitempty"`
}

type ContactId struct {
	EmailAddress string `xml:"EmailAddress,omitempty" json:"EmailAddress,omitempty"`
	PhoneNumber  string `xml:"PhoneNumber,omitempty" json:"PhoneNumber,omitempty"`
	FaxNumber    string `xml:"FaxNumber,omitempty" json:"FaxNumber,omitempty"`
}

type GoverningAgreementType struct {
	Value            string `xml:",chardata"`
	Namespace        string `xml:"Namespace,attr,omitempty" json:"Namespace,omitempty"`
	UserDefinedValue string `xml:"UserDefinedValue,attr,omitempty" json:"UserDefinedValue,omitempty"`
}

type Performance struct {
	Territory *TerritoryCode `xml:"Territory,omitempty" json:"Territory,omitempty"`
	Date      *EventDate     `xml:"EventDate,omitempty" json:"EventDate,omitempty"`
}

type ArtistRole struct {
	Value            string `xml:",chardata"`
	Namespace        string `xml:"Namespace,attr,omitempty" json:"Namespace,omitempty"`
	UserDefinedValue string `xml:"UserDefinedValue,attr,omitempty" json:"UserDefinedValue,omitempty"`
}

type ArtistDelegatedUsageRights struct {
	UseType                     []*UseType         `xml:"UseType,omitempty" json:"UseType,omitempty"`
	UserInterfaceType           *UserInterfaceType `xml:"UserInterfaceType,omitempty" json:"UserInterfaceType,omitempty"`
	PeriodOfRightsDelegation    *Period            `xml:"PeriodOfRightsDelegation,omitempty" json:"PeriodOfRightsDelegation,omitempty"`
	TerritoryOfRightsDelegation *TerritoryCode     `xml:"TerritoryOfRightsDelegation,omitempty" json:"TerritoryOfRightsDelegation,omitempty"`
	MembershipType              string             `xml:"MembershipType,omitempty" json:"MembershipType,omitempty"`
}

type CreationId struct {
	ISWC                  string         `xml:"ISWC,omitempty" json:"ISWC,omitempty"`
	OpusNumber            string         `xml:"OpusNumber,omitempty" json:"OpusNumber,omitempty"`
	ComposerCatalogNumber string         `xml:"ComposerCatalogNumber,omitempty" json:"ComposerCatalogNumber,omitempty"`
	ISRC                  string         `xml:"ISRC,omitempty" json:"ISRC,omitempty"`
	ISMN                  string         `xml:"ISMN,omitempty" json:"ISMN,omitempty"`
	ISAN                  string         `xml:"ISAN,omitempty" json:"ISAN,omitempty"`
	VISAN                 string         `xml:"VISAN,omitempty" json:"VISAN,omitempty"`
	ISBN                  string         `xml:"ISBN,omitempty" json:"ISBN,omitempty"`
	ISSN                  string         `xml:"ISSN,omitempty" json:"ISSN,omitempty"`
	SICI                  string         `xml:"SICI,omitempty" json:"SICI,omitempty"`
	CatalogNumber         *CatalogNumber `xml:"CatalogNumber,omitempty" json:"CatalogNumber,omitempty"`
	ProprietaryId         *ProprietaryId `xml:"ProprietaryId,omitempty" json:"ProprietaryId,omitempty"`
}

type CatalogNumber struct {
	Value     string `xml:",chardata"`
	Namespace string `xml:"Namespace,attr,omitempty" json:"Namespace,omitempty"`
}

type CueUseType struct {
	Value            string `xml:",chardata"`
	Namespace        string `xml:"Namespace,attr,omitempty" json:"Namespace,omitempty"`
	UserDefinedValue string `xml:"UserDefinedValue,attr,omitempty" json:"UserDefinedValue,omitempty"`
}

type CueThemeType struct {
	Value            string `xml:",chardata"`
	Namespace        string `xml:"Namespace,attr,omitempty" json:"Namespace,omitempty"`
	UserDefinedValue string `xml:"UserDefinedValue,attr,omitempty" json:"UserDefinedValue,omitempty"`
}

type CueVocalType struct {
	Value            string `xml:",chardata"`
	Namespace        string `xml:"Namespace,attr,omitempty" json:"Namespace,omitempty"`
	UserDefinedValue string `xml:"UserDefinedValue,attr,omitempty" json:"UserDefinedValue,omitempty"`
}

type CueVisualPerceptionType struct {
	Value            string `xml:",chardata"`
	Namespace        string `xml:"Namespace,attr,omitempty" json:"Namespace,omitempty"`
	UserDefinedValue string `xml:"UserDefinedValue,attr,omitempty" json:"UserDefinedValue,omitempty"`
}

type CueOrigin struct {
	Value            string `xml:",chardata"`
	Namespace        string `xml:"Namespace,attr,omitempty" json:"Namespace,omitempty"`
	UserDefinedValue string `xml:"UserDefinedValue,attr,omitempty" json:"UserDefinedValue,omitempty"`
}

type CueCreationReference struct {
	CueWorkReference     string `xml:"CueWorkReference,omitempty" json:"CueWorkReference,omitempty"`
	CueResourceReference string `xml:"CueResourceReference,omitempty" json:"CueResourceReference,omitempty"`
}

type CueSheetType struct {
	Value            string `xml:",chardata"`
	Namespace        string `xml:"Namespace,attr,omitempty" json:"Namespace,omitempty"`
	UserDefinedValue string `xml:"UserDefinedValue,attr,omitempty" json:"UserDefinedValue,omitempty"`
}

type WorkList struct {
	MusicalWork           []*MusicalWork `xml:"MusicalWork,omitempty" json:"MusicalWork,omitempty"`
	LanguageAndScriptCode string         `xml:"LanguageAndScriptCode,attr,omitempty" json:"LanguageAndScriptCode,omitempty"`
}

type MusicalWork struct {
	MusicalWorkId                 *MusicalWorkId                   `xml:"MusicalWorkId,omitempty" json:"MusicalWorkId,omitempty"`
	MusicalWorkReference          string                           `xml:"MusicalWorkReference,omitempty" json:"MusicalWorkReference,omitempty"`
	ReferenceTitle                *ReferenceTitle                  `xml:"ReferenceTitle,omitempty" json:"ReferenceTitle,omitempty"`
	RightsAgreementId             *RightsAgreementId               `xml:"RightsAgreementId,omitempty" json:"RightsAgreementId,omitempty"`
	MusicalWorkContributor        []*MusicalWorkContributor        `xml:"MusicalWorkContributor,omitempty" json:"MusicalWorkContributor,omitempty"`
	MusicalWorkType               *MusicalWorkType                 `xml:"MusicalWorkType,omitempty" json:"MusicalWorkType,omitempty"`
	RightShare                    *RightShare                      `xml:"RightShare,omitempty" json:"RightShare,omitempty"`
	MusicalWorkDetailsByTerritory []*MusicalWorkDetailsByTerritory `xml:"MusicalWorkDetailsByTerritory,omitempty" json:"MusicalWorkDetailsByTerritory,omitempty"`
	IsUpdated                     bool                             `xml:"IsUpdated,attr,omitempty" json:"IsUpdated,omitempty"`
	LanguageAndScriptCode         string                           `xml:"LanguageAndScriptCode,attr,omitempty" json:"LanguageAndScriptCode,omitempty"`
}

type MusicalWorkId struct {
	ISWC                  string         `xml:"ISWC,omitempty" json:"ISWC,omitempty"`
	OpusNumber            string         `xml:"OpusNumber,omitempty" json:"OpusNumber,omitempty"`
	ComposerCatalogNumber string         `xml:"ComposerCatalogNumber,omitempty" json:"ComposerCatalogNumber,omitempty"`
	ProprietaryId         *ProprietaryId `xml:"ProprietaryId,omitempty" json:"ProprietaryId,omitempty"`
	IsReplaced            bool           `xml:"IsReplaced,attr,omitempty" json:"IsReplaced,omitempty"`
}

type ReferenceTitle struct {
	TitleText             *TitleText `xml:"TitleText,omitempty" json:"TitleText,omitempty"`
	SubTitle              *SubTitle  `xml:"SubTitle,omitempty" json:"SubTitle,omitempty"`
	LanguageAndScriptCode string     `xml:"LanguageAndScriptCode,attr,omitempty" json:"LanguageAndScriptCode,omitempty"`
}

type TitleText struct {
	Value                 string `xml:",chardata"`
	LanguageAndScriptCode string `xml:"LanguageAndScriptCode,attr,omitempty" json:"LanguageAndScriptCode,omitempty"`
}

type SubTitle struct {
	Value                 string `xml:",chardata"`
	LanguageAndScriptCode string `xml:"LanguageAndScriptCode,attr,omitempty" json:"LanguageAndScriptCode,omitempty"`
}

type ProprietaryId struct {
	Value     string `xml:",chardata"`
	Namespace string `xml:"Namespace,attr,omitempty" json:"Namespace,omitempty"`
}

type MusicalWorkDetailsByTerritory struct {
	TerritoryCode          *TerritoryCode            `xml:"TerritoryCode,omitempty" json:"TerritoryCode,omitempty"`
	ExcludedTerritoryCode  *TerritoryCode            `xml:"ExcludedTerritoryCode,omitempty" json:"ExcludedTerritoryCode,omitempty"`
	MusicalWorkContributor []*MusicalWorkContributor `xml:"MusicalWorkContributor,omitempty" json:"MusicalWorkContributor,omitempty"`
	DisplayArtistName      *Name                     `xml:"DisplayArtistName,omitempty" json:"DisplayArtistName,omitempty"`
	LanguageAndScriptCode  string                    `xml:"LanguageAndScriptCode,attr,omitempty" json:"LanguageAndScriptCode,omitempty"`
}

type RightsAgreementId struct {
	MWLI          string         `xml:"MWLI,omitempty" json:"MWLI,omitempty"`
	ProprietaryId *ProprietaryId `xml:"ProprietaryId,omitempty" json:"ProprietaryId,omitempty"`
}

type MusicalWorkContributor struct {
	PartyId                    *PartyId                    `xml:"PartyId,omitempty" json:"PartyId,omitempty"`
	PartyName                  *PartyName                  `xml:"PartyName,omitempty" json:"PartyName,omitempty"`
	MusicalWorkContributorRole *MusicalWorkContributorRole `xml:"MusicalWorkContributorRole,omitempty" json:"MusicalWorkContributorRole,omitempty"`
	SocietyAffiliation         *SocietyAffiliation         `xml:"SocietyAffiliation,omitempty" json:"SocietyAffiliation,omitempty"`
	SequenceNumber             int                         `xml:"SequenceNumber,attr,omitempty" json:"SequenceNumber,omitempty"`
}

type SocietyAffiliation struct {
	TerritoryCode         *TerritoryCode   `xml:"TerritoryCode,omitempty" json:"TerritoryCode,omitempty"`
	ExcludedTerritoryCode *TerritoryCode   `xml:"ExcludedTerritoryCode,omitempty" json:"ExcludedTerritoryCode,omitempty"`
	MusicRightsSociety    *PartyDescriptor `xml:"MusicRightsSociety,omitempty" json:"MusicRightsSociety,omitempty"`
}

type MusicalWorkContributorRole struct {
	Value            string `xml:",chardata"`
	Namespace        string `xml:"Namespace,attr,omitempty" json:"Namespace,omitempty"`
	UserDefinedValue string `xml:"UserDefinedValue,attr,omitempty" json:"UserDefinedValue,omitempty"`
}

type MusicalWorkType struct {
	Value            string `xml:",chardata"`
	Namespace        string `xml:"Namespace,attr,omitempty" json:"Namespace,omitempty"`
	UserDefinedValue string `xml:"UserDefinedValue,attr,omitempty" json:"UserDefinedValue,omitempty"`
}

type RightShare struct {
	RightShareId                    *RightsAgreementId               `xml:"RightShareId,omitempty" json:"RightShareId,omitempty"`
	RightShareReference             string                           `xml:"RightShareReference,omitempty" json:"RightShareReference,omitempty"`
	RightShareCreationReferenceList *RightShareCreationReferenceList `xml:"RightShareCreationReferenceList,omitempty" json:"RightShareCreationReferenceList,omitempty"`
	TerritoryCode                   *TerritoryCode                   `xml:"TerritoryCode,omitempty" json:"TerritoryCode,omitempty"`
	ExcludedTerritoryCode           *TerritoryCode                   `xml:"ExcludedTerritoryCode,omitempty" json:"ExcludedTerritoryCode,omitempty"`
	RightsType                      *RightsType                      `xml:"RightsType,omitempty" json:"RightsType,omitempty"`
	UseType                         []*UseType                       `xml:"UseType,omitempty" json:"UseType,omitempty"`
	UserInterfaceType               *UserInterfaceType               `xml:"UserInterfaceType,omitempty" json:"UserInterfaceType,omitempty"`
	DistributionChannelType         *DistributionChannelType         `xml:"DistributionChannelType,omitempty" json:"DistributionChannelType,omitempty"`
	CarrierType                     *CarrierType                     `xml:"CarrierType,omitempty" json:"CarrierType,omitempty"`
	CommercialModelType             *CommercialModelType             `xml:"CommercialModelType,omitempty" json:"CommercialModelType,omitempty"`
	MusicalWorkRightsClaimType      string                           `xml:"MusicalWorkRightsClaimType,omitempty" json:"MusicalWorkRightsClaimType,omitempty"`
	RightsController                []*RightsController              `xml:"RightsController,omitempty" json:"RightsController,omitempty"`
	ValidityPeriod                  *Period                          `xml:"ValidityPeriod,omitempty" json:"ValidityPeriod,omitempty"`
	RightShareUnknown               bool                             `xml:"RightShareUnknown,omitempty" json:"RightShareUnknown,omitempty"`
	RightSharePercentage            *Percentage                      `xml:"RightSharePercentage,omitempty" json:"RightSharePercentage,omitempty"`
	TariffReference                 *TariffReference                 `xml:"TariffReference,omitempty" json:"TariffReference,omitempty"`
	LicenseStatus                   string                           `xml:"LicenseStatus,omitempty" json:"LicenseStatus,omitempty"`
	HasFirstLicenseRefusal          bool                             `xml:"HasFirstLicenseRefusal,omitempty" json:"HasFirstLicenseRefusal,omitempty"`
	LanguageAndScriptCode           string                           `xml:"LanguageAndScriptCode,attr,omitempty" json:"LanguageAndScriptCode,omitempty"`
}

type RightShareCreationReferenceList struct {
	RightShareWorkReference     []string `xml:"RightShareWorkReference,omitempty" json:"RightShareWorkReference,omitempty"`
	RightShareResourceReference []string `xml:"RightShareResourceReference,omitempty" json:"RightShareResourceReference,omitempty"`
	RightShareReleaseReference  []string `xml:"RightShareReleaseReference,omitempty" json:"RightShareReleaseReference,omitempty"`
}

type RightsType struct {
	Value            string `xml:",chardata"`
	TerritoryCode    string `xml:"TerritoryCode,attr,omitempty" json:"TerritoryCode,omitempty"`
	Namespace        string `xml:"Namespace,attr,omitempty" json:"Namespace,omitempty"`
	UserDefinedValue string `xml:"UserDefinedValue,attr,omitempty" json:"UserDefinedValue,omitempty"`
}

type UseType struct {
	Value            string `xml:",chardata"`
	Namespace        string `xml:"Namespace,attr,omitempty" json:"Namespace,omitempty"`
	UserDefinedValue string `xml:"UserDefinedValue,attr,omitempty" json:"UserDefinedValue,omitempty"`
}

type UserInterfaceType struct {
	Value            string `xml:",chardata"`
	Namespace        string `xml:"Namespace,attr,omitempty" json:"Namespace,omitempty"`
	UserDefinedValue string `xml:"UserDefinedValue,attr,omitempty" json:"UserDefinedValue,omitempty"`
}

type DistributionChannelType struct {
	Value            string `xml:",chardata"`
	Namespace        string `xml:"Namespace,attr,omitempty" json:"Namespace,omitempty"`
	UserDefinedValue string `xml:"UserDefinedValue,attr,omitempty" json:"UserDefinedValue,omitempty"`
}

type CarrierType struct {
	Value            string `xml:",chardata"`
	Namespace        string `xml:"Namespace,attr,omitempty" json:"Namespace,omitempty"`
	UserDefinedValue string `xml:"UserDefinedValue,attr,omitempty" json:"UserDefinedValue,omitempty"`
}

type CommercialModelType struct {
	Value            string `xml:",chardata"`
	Namespace        string `xml:"Namespace,attr,omitempty" json:"Namespace,omitempty"`
	UserDefinedValue string `xml:"UserDefinedValue,attr,omitempty" json:"UserDefinedValue,omitempty"`
}

type RightsController struct {
	PartyId              *PartyId    `xml:"PartyId,omitempty" json:"PartyId,omitempty"`
	PartyName            *PartyName  `xml:"PartyName,omitempty" json:"PartyName,omitempty"`
	RightsControllerRole string      `xml:"RightsControllerRole,omitempty" json:"RightsControllerRole,omitempty"`
	RightShareUnknown    bool        `xml:"RightShareUnknown,omitempty" json:"RightShareUnknown,omitempty"`
	RightSharePercentage *Percentage `xml:"RightSharePercentage,omitempty" json:"RightSharePercentage,omitempty"`
	RightsControllerType string      `xml:"RightsControllerType,omitempty" json:"RightsControllerType,omitempty"`
	SequenceNumber       int         `xml:"SequenceNumber,attr,omitempty" json:"SequenceNumber,omitempty"`
}

type TypedRightsController struct {
	PartyId                 *PartyId       `xml:"PartyId,omitempty" json:"PartyId,omitempty"`
	PartyName               *PartyName     `xml:"PartyName,omitempty" json:"PartyName,omitempty"`
	RightsControllerRole    string         `xml:"RightsControllerRole,omitempty" json:"RightsControllerRole,omitempty"`
	RightShareUnknown       bool           `xml:"RightShareUnknown,omitempty" json:"RightShareUnknown,omitempty"`
	RightSharePercentage    *Percentage    `xml:"RightSharePercentage,omitempty" json:"RightSharePercentage,omitempty"`
	RightsControllerType    string         `xml:"RightsControllerType,omitempty" json:"RightsControllerType,omitempty"`
	TerritoryOfRegistration *TerritoryCode `xml:"TerritoryOfRegistration,omitempty" json:"TerritoryOfRegistration,omitempty"`
	StartDate               string         `xml:"StartDate,omitempty" json:"StartDate,omitempty"`
	EndDate                 string         `xml:"EndDate,omitempty" json:"EndDate,omitempty"`
	SequenceNumber          int            `xml:"SequenceNumber,attr,omitempty" json:"SequenceNumber,omitempty"`
}

type TerritoryCode struct {
	Value          string `xml:",chardata"`
	IdentifierType string `xml:"IdentifierType,attr,omitempty" json:"IdentifierType,omitempty"`
}

type Period struct {
	StartDate     *EventDate     `xml:"StartDate,omitempty" json:"StartDate,omitempty"`
	EndDate       *EventDate     `xml:"EndDate,omitempty" json:"EndDate,omitempty"`
	StartDateTime *EventDateTime `xml:"StartDateTime,omitempty" json:"StartDateTime,omitempty"`
	EndDateTime   *EventDateTime `xml:"EndDateTime,omitempty" json:"EndDateTime,omitempty"`
}

type EventDate struct {
	Value                 string `xml:",chardata"`
	IsApproximate         bool   `xml:"IsApproximate,attr,omitempty" json:"IsApproximate,omitempty"`
	IsBefore              bool   `xml:"IsBefore,attr,omitempty" json:"IsBefore,omitempty"`
	IsAfter               bool   `xml:"IsAfter,attr,omitempty" json:"IsAfter,omitempty"`
	TerritoryCode         string `xml:"TerritoryCode,attr,omitempty" json:"TerritoryCode,omitempty"`
	LocationDescription   string `xml:"LocationDescription,attr,omitempty" json:"LocationDescription,omitempty"`
	LanguageAndScriptCode string `xml:"LanguageAndScriptCode,attr,omitempty" json:"LanguageAndScriptCode,omitempty"`
}

type EventDateTime struct {
	Value                 string `xml:",chardata"`
	IsApproximate         bool   `xml:"IsApproximate,attr,omitempty" json:"IsApproximate,omitempty"`
	IsBefore              bool   `xml:"IsBefore,attr,omitempty" json:"IsBefore,omitempty"`
	IsAfter               bool   `xml:"IsAfter,attr,omitempty" json:"IsAfter,omitempty"`
	TerritoryCode         string `xml:"TerritoryCode,attr,omitempty" json:"TerritoryCode,omitempty"`
	LocationDescription   string `xml:"LocationDescription,attr,omitempty" json:"LocationDescription,omitempty"`
	LanguageAndScriptCode string `xml:"LanguageAndScriptCode,attr,omitempty" json:"LanguageAndScriptCode,omitempty"`
}

type Percentage struct {
	Value            string `xml:",chardata"`
	HasMaxValueOfOne bool   `xml:"HasMaxValueOfOne,attr,omitempty" json:"HasMaxValueOfOne,omitempty"`
}

type TariffReference struct {
	Value                 string `xml:",chardata"`
	LanguageAndScriptCode string `xml:"LanguageAndScriptCode,attr,omitempty" json:"LanguageAndScriptCode,omitempty"`
	TariffSubReference    string `xml:"TariffSubReference,attr,omitempty" json:"TariffSubReference,omitempty"`
}
