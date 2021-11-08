import { MatIconRegistry } from "@angular/material/icon";
import { DomSanitizer } from "@angular/platform-browser";

export function RegisterCustomIcons(matIconRegistry: MatIconRegistry, domSanitizer: DomSanitizer) {
  matIconRegistry.addSvgIcon(
    'github',
    domSanitizer.bypassSecurityTrustResourceUrl('http://localhost:8080/static/frontend/dist/assets/icons/github.svg')
  )

  matIconRegistry.addSvgIcon(
    'qrcode_scan',
    domSanitizer.bypassSecurityTrustResourceUrl('http://localhost:8080/static/frontend/dist/assets/icons/qrcode-scan.svg')
  )
}
