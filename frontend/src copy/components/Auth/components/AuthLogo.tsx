import React from 'react';
import { AuthLogoContainer, AuthLogoImg } from '../styles.styled';
import logo from '../../../assets/imgs/logo.jpg';
type IProps = {};
const AuthLogo: React.FC<IProps> = () => {
  return (
    <AuthLogoContainer>
      <AuthLogoImg src={logo} />
    </AuthLogoContainer>
  );
};

export default AuthLogo;
