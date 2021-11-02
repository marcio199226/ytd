import { APP_BASE_HREF, CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { BrowserModule } from '@angular/platform-browser';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { HttpClient, HttpClientModule } from '@angular/common/http';

import { AppRoutingModule } from './app-routing.module';
import { AppComponent } from './app.component';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';
import { HomeComponent, OfflinePlaylistComponent } from './pages';
import { FlexLayoutModule } from '@angular/flex-layout';
import { MatButtonModule } from '@angular/material/button';
import { MatBadgeModule } from '@angular/material/badge';
import { MatCardModule } from '@angular/material/card';
import { MatCheckboxModule } from '@angular/material/checkbox';
import { MatDialogModule, MAT_DIALOG_DEFAULT_OPTIONS } from '@angular/material/dialog';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatIconModule } from '@angular/material/icon';
import { MatInputModule } from '@angular/material/input';
import { MatMenuModule } from '@angular/material/menu';
import { MatProgressBarModule } from '@angular/material/progress-bar';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MatSelectModule } from '@angular/material/select';
import { MatSlideToggleModule } from '@angular/material/slide-toggle';
import { MatSliderModule}  from '@angular/material/slider';
import { MatSidenavModule } from '@angular/material/sidenav';
import { MatSnackBarModule, MAT_SNACK_BAR_DEFAULT_OPTIONS } from '@angular/material/snack-bar';
import { MatTabsModule } from '@angular/material/tabs';
import { MatToolbarModule } from '@angular/material/toolbar';
import { MatTooltipModule } from '@angular/material/tooltip';
import { NG_EVENT_PLUGINS } from '@tinkoff/ng-event-plugins';
import { TranslateLoader, TranslateModule } from '@ngx-translate/core';
import { TranslateHttpLoader } from '@ngx-translate/http-loader';
import { QrCodeModule } from 'ng-qrcode';
import {
  AudioPlayerComponent,
  LoaderOverlayComponent,
  SettingsComponent,
  UpdaterComponent,
  ConfirmationDialogComponent,
  AddToPlaylistComponent,
  CreatePlaylistComponent,
  LanguageSwitcherComponent
} from './components';
import { AutofocusDirective } from './directives';
import { environment } from 'environments/environment';

export function HttpLoaderFactory(http: HttpClient) {
  const assetsDir = environment.production ? 'http://localhost:8080/static/frontend/dist/assets/i18n/' : './assets/i18n/'
  return new TranslateHttpLoader(http, assetsDir, '.json');
}

@NgModule({
  declarations: [
    AppComponent,
    HomeComponent,
    OfflinePlaylistComponent,
    AudioPlayerComponent,
    SettingsComponent,
    UpdaterComponent,
    AddToPlaylistComponent,
    CreatePlaylistComponent,
    LanguageSwitcherComponent,
    ConfirmationDialogComponent,
    LoaderOverlayComponent,
    // directives
    AutofocusDirective
  ],
  imports: [
    BrowserModule,
    CommonModule,
    HttpClientModule,
    AppRoutingModule,
    BrowserAnimationsModule,
    FormsModule,
    ReactiveFormsModule,
    TranslateModule.forRoot({
      defaultLanguage: 'en',
      loader: {
        provide: TranslateLoader,
        useFactory: HttpLoaderFactory,
        deps: [HttpClient]
      }
    }),
    FlexLayoutModule,
    QrCodeModule,
    MatBadgeModule,
    MatButtonModule,
    MatCardModule,
    MatCheckboxModule,
    MatDialogModule,
    MatFormFieldModule,
    MatIconModule,
    MatInputModule,
    MatMenuModule,
    MatProgressBarModule,
    MatProgressSpinnerModule,
    MatSelectModule,
    MatSlideToggleModule,
    MatSliderModule,
    MatSidenavModule,
    MatSnackBarModule,
    MatTabsModule,
    MatToolbarModule,
    MatTooltipModule
  ],
  providers: [
    {provide: APP_BASE_HREF, useValue : '/' },
    NG_EVENT_PLUGINS,
    { provide: MAT_SNACK_BAR_DEFAULT_OPTIONS, useValue: { duration: 5000 } },
    {
      provide: MAT_DIALOG_DEFAULT_OPTIONS,
      useValue: {
        backdropClass: 'blurred-backdrop-bg',
        hasBackdrop: true,
        autoFocus: false
      }
    },
  ],
  bootstrap: [AppComponent]
})
export class AppModule { }
