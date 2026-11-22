import './css/tema.css';
import './css/jacobitus.css';
import './css/pesaje.css';
import './style.css';
import './app.css';

import { montarLogin } from './login';

const app = document.querySelector('#app') as HTMLElement;
montarLogin(app);
