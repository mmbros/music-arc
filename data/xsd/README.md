### Validate Schema XSD

How to validate the `music-arc.xsd`:

    xmllint --schema http://www.w3.org/2001/XMLSchema.xsd --noout music-arc.xsd

### Validate XML against Schema XSD
How to check if `music-arc-inc.xml` is a valid document against the `music-arc.xsd` schema definition:

    xmllint --schema music-arc.xsd --noout --noent music-arc-inc.xml

--noent serve nel caso in cui siano definite ENTITA'

<!DOCTYPE music-arc [
  <!ENTITY sep    "&#x2022;"><!-- Bullet -->
	<!ENTITY bull   "&#x2022;"><!-- Bullet -->
  <!ENTITY hellip "&#x2026;"><!-- Horizontal ellipsis -->
  <!ENTITY nbsp     "&#xA0;"><!-- Non-breaking space -->
]>
