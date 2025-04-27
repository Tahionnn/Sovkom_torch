import { createUrl } from '../../other';
import axios from 'axios';
import { apies } from '../../../lists/apies';
class AuthApi {
  async register(username: string, password: string) {
    try {
      if (username && password) {
        const user = await axios.post(
          createUrl(apies.auth.authHost, [apies.auth.register]),
          { username, password }
        );
        localStorage.setItem('token', user.data['token']);
        return user;
      }
    } catch (error) {
      throw error;
    }
  }

  async login(username: string, password: string) {
    try {
      if (username && password) {
        const user: any = await axios.post(
          createUrl(apies.auth.authHost, [apies.auth.login]),
          { username, password }
        );
        localStorage.setItem('token', user.data['token']);
        return user;
      }
    } catch (error) {
      throw error;
    }
  }
}

export default new AuthApi();
