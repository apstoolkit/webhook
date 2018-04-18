package main

import (
	"encoding/xml"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

var samplePayload = `
<?xml version="1.0" encoding="utf-8"?>
<DocuSignEnvelopeInformation
    xmlns:xsd="http://www.w3.org/2001/XMLSchema"
    xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
    xmlns="http://www.docusign.net/API/3.0">
    <EnvelopeStatus>
        <RecipientStatuses>
            <RecipientStatus>
                <Type>Signer</Type>
                <Email>user@domain</Email>
                <UserName>Bob Boone</UserName>
                <RoutingOrder>1</RoutingOrder>
                <Sent>2018-04-17T14:59:28.92</Sent>
                <Delivered>2018-04-17T15:00:01.887</Delivered>
                <Signed>2018-04-17T15:00:08.183</Signed>
                <DeclineReason xsi:nil="true" />
                <Status>Completed</Status>
                <RecipientIPAddress>111.22.33.44</RecipientIPAddress>
                <CustomFields />
                <TabStatuses>
                    <TabStatus>
                        <TabType>SignHere</TabType>
                        <Status>Signed</Status>
                        <XPosition>295</XPosition>
                        <YPosition>245</YPosition>
                        <TabLabel>signer1sig</TabLabel>
                        <TabName>Please sign here</TabName>
                        <TabValue />
                        <DocumentID>1</DocumentID>
                        <PageNumber>1</PageNumber>
                    </TabStatus>
                    <TabStatus>
                        <TabType>Custom</TabType>
                        <Status>Signed</Status>
                        <XPosition>268</XPosition>
                        <YPosition>518</YPosition>
                        <TabLabel>City</TabLabel>
                        <TabName>City</TabName>
                        <TabValue>London</TabValue>
                        <DocumentID>1</DocumentID>
                        <PageNumber>1</PageNumber>
                        <OriginalValue>London</OriginalValue>
                        <CustomTabType>Text</CustomTabType>
                    </TabStatus>
                    <TabStatus>
                        <TabType>FullName</TabType>
                        <Status>Signed</Status>
                        <XPosition>293</XPosition>
                        <YPosition>439</YPosition>
                        <TabLabel>Full Name</TabLabel>
                        <TabName>Full Name</TabName>
                        <TabValue>Doug Smith</TabValue>
                        <DocumentID>1</DocumentID>
                        <PageNumber>1</PageNumber>
                    </TabStatus>
                </TabStatuses>
                <RecipientAttachment>
                    <Attachment>
                        <Data>xxxx</Data>
                        <Label>DSXForm</Label>
                    </Attachment>
                </RecipientAttachment>
                <AccountStatus>Active</AccountStatus>
                <FormData>
                    <xfdf>
                        <fields>
                            <field name="City">
                                <value>London</value>
                            </field>
                            <field name="FullName">
                                <value>Bob Boone</value>
                            </field>
                        </fields>
                    </xfdf>
                </FormData>
                <RecipientId>5dc2715e-2efc-4c2b-ab4b-208a0c030c0e</RecipientId>
            </RecipientStatus>
        </RecipientStatuses>
        <TimeGenerated>2018-04-17T15:00:22.2587658</TimeGenerated>
        <EnvelopeID>dd8449dc-3989-4b81-8d72-5e927a5769ad</EnvelopeID>
        <Subject>NewCo agreement for signature</Subject>
        <UserName>Docusign Guy</UserName>
        <Email>docusign_acct@somedomain</Email>
        <Status>Completed</Status>
        <Created>2018-04-17T14:59:28.277</Created>
        <Sent>2018-04-17T14:59:28.95</Sent>
        <Delivered>2018-04-17T15:00:02.013</Delivered>
        <Signed>2018-04-17T15:00:08.183</Signed>
        <Completed>2018-04-17T15:00:08.183</Completed>
        <ACStatus>Original</ACStatus>
        <ACStatusDate>2018-04-17T14:59:28.277</ACStatusDate>
        <ACHolder>Docusign Guy</ACHolder>
        <ACHolderEmail>docusign_acct@somedomain</ACHolderEmail>
        <ACHolderLocation>DocuSign</ACHolderLocation>
        <SigningLocation>Online</SigningLocation>
        <SenderIPAddress>111.22.33.44</SenderIPAddress>
        <EnvelopePDFHash />
        <CustomFields>
            <CustomField>
                <Name>AccountId</Name>
                <Show>false</Show>
                <Required>false</Required>
                <Value>111</Value>
                <CustomFieldType>Text</CustomFieldType>
            </CustomField>
            <CustomField>
                <Name>AccountName</Name>
                <Show>false</Show>
                <Required>false</Required>
                <Value>DougSmith</Value>
                <CustomFieldType>Text</CustomFieldType>
            </CustomField>
            <CustomField>
                <Name>AccountSite</Name>
                <Show>false</Show>
                <Required>false</Required>
                <Value>demo</Value>
                <CustomFieldType>Text</CustomFieldType>
            </CustomField>
        </CustomFields>
        <AutoNavigation>true</AutoNavigation>
        <EnvelopeIdStamping>true</EnvelopeIdStamping>
        <AuthoritativeCopy>false</AuthoritativeCopy>
        <DocumentStatuses>
            <DocumentStatus>
                <ID>1</ID>
                <Name>Agreement &amp; test</Name>
                <TemplateName />
                <Sequence>1</Sequence>
            </DocumentStatus>
        </DocumentStatuses>
    </EnvelopeStatus>
    <DocumentPDFs>
        <DocumentPDF>
            <Name>CertificateOfCompletion_dd8449dc-3989-4b81-8d72-5e927a5769ad.pdf</Name>
            <PDFBytes>xxx</PDFBytes>
            <DocumentType>SUMMARY</DocumentType>
        </DocumentPDF>
    </DocumentPDFs>
    <TimeZone>Pacific Standard Time</TimeZone>
    <TimeZoneOffset>-7</TimeZoneOffset>
</DocuSignEnvelopeInformation>`

func TestParsing(t *testing.T) {
	var envInfo DocuSignEnvelopeInformation

	err := xml.Unmarshal([]byte(samplePayload), &envInfo)
	if assert.Nil(t, err) {
		fmt.Printf("%+v", envInfo)
		assert.Equal(t, "2018-04-17T15:00:22.2587658", envInfo.EnvelopeStatus.TimeGenerated)
	}
}
