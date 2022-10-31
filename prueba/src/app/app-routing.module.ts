import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { InicioComponent } from './component/inicio/inicio.component'
import { LoginComponent } from './component/login/login.component'
import { ReportesComponent } from './component/reportes/reportes.component';
import { Reportes2Component } from './component/reportes2/reportes2.component';
import { Reportes3Component } from './component/reportes3/reportes3.component';

const routes: Routes = [

  { path: 'inicio', component: InicioComponent},
  { path: 'login', component: LoginComponent},
  { path: 'reporte', component: ReportesComponent},
  { path: 'reporte2', component: Reportes2Component},
  { path: 'reporte3', component: Reportes3Component},
  { path: '**', redirectTo: 'inicio' }

];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})
export class AppRoutingModule { }
