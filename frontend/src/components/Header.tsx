import { HeaderContainer, HeaderLogoImg } from './Header.styled';
import logo from '../assets/imgs/logo.jpg';
type Props = {};
const Header: React.FC<Props> = () => {
  return (
    <HeaderContainer>
      <HeaderLogoImg src={logo}></HeaderLogoImg>
    </HeaderContainer>
  );
};

export default Header;
