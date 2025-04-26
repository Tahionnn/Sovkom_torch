import { TextField } from '@mui/material';
import { useEffect, useState } from 'react';
import {
  AuthButton,
  AuthCard,
  AuthContainer,
  AuthTitle,
} from './styles.styled';
import React from 'react';
import AuthLogo from './components/AuthLogo';
import Auth from '../../utils/requests/auth/AuthApi';
import { createCookie, useNavigate } from 'react-router';
import { routes } from '../../lists/routes';
import Header from '../Header';
type Props = {};
const AuthPage: React.FC<Props> = () => {
  const [checkAuthVariant, setCheckAuthVariant] = useState<boolean>(true);
  const [username, setUsername] = useState<string>('');
  const [password, setPassword] = useState<string>('');
  const [checkAuth, setCheckAuth] = useState<boolean>(false);
  const navigate = useNavigate();
  useEffect(() => {
    if (checkAuth) {
      navigate(routes.start);
    }
  }, [checkAuth]);
  return (
    <>
      <AuthContainer>
        <AuthCard>
          <AuthLogo />
          {!checkAuthVariant ? (
            <AuthTitle
              onClick={() => {
                setCheckAuthVariant(true);
              }}
            >
              {' '}
              Регистрация
            </AuthTitle>
          ) : (
            <AuthTitle
              onClick={() => {
                setCheckAuthVariant(false);
              }}
            >
              Вход
            </AuthTitle>
          )}
          <TextField
            placeholder="Имя пользователя"
            value={username}
            onChange={(e) => {
              setUsername(e.target.value);
            }}
          ></TextField>
          <TextField
            placeholder="Пароль"
            value={password}
            type="password"
            onChange={(e) => {
              setPassword(e.target.value);
            }}
          ></TextField>
          {checkAuthVariant ? (
            <AuthButton
              onClick={() => {
                Auth.register(username, password)
                  .then((res) => {
                    createCookie('_token', res?.data.token);
                    setCheckAuth(true);
                  })
                  .catch((err) => {
                    setCheckAuth(false);
                  });
              }}
            >
              Зарегистрироваться
            </AuthButton>
          ) : (
            <AuthButton
              onClick={() => {
                Auth.login(username, password)
                  .then((res) => {
                    setCheckAuth(true);
                  })
                  .catch((err) => {
                    setCheckAuth(false);
                  });
              }}
            >
              Войти
            </AuthButton>
          )}
        </AuthCard>
      </AuthContainer>
    </>
  );
};

export default AuthPage;
