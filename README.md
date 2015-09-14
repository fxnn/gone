# gone

This is _Work In Progress_, the following information is only a raw sketch / a collection of ideas.

## Ziele

* KISS-Wiki
* das sich im Dateisystem widerspiegelt
* mit dem auch Konfigurationsdateien angezeigt werden können
* mit dem auch vorhandene Dokumentation wie Manpages angezeigt werden kann


## Konzepte

### Abbildung URL-Dateisystem

* URLs werden ins Dateisystem gematcht
* Konfigurationsdatei regelt Abbildungen von URLs in Verzeichnisse
* Konfigurierbare Input-Pipeline und Output-Pipeline je Dateityp, d.h.
  * ich kann etwa einen Markdown-HTML-Parser installieren,
  * einen Text-HTML-Renderer oder auch
  * einen Manpage-HTML-Renderer.
  * Die Output-Pipeline regelt, wie der vom User eingegebene Text in die Datei umgewandelt wird.
  * Ohne Konfiguration keine Möglichkeit, den Dateitypen zu speichern/lesen.
* Query-Parameter erlauben in einen Schreib- oder Info-Modus zu wechseln
* Später könnte es eine History automatisch für git-versionierte Verzeichnisse geben

### Authentifizierung

* HTTP-Basic-Auth <-> htpasswd?

