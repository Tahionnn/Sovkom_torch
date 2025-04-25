import { TextField } from '@mui/material';
import React, { useEffect } from "react"
import { AuthButton, AuthCard, AuthContainer, AuthTitle } from './styles.styled.ts';
import MainStore from '../../stores/Main.store.ts';
type Props={}

const AuthPage: React.FC<Props> = () => {
    return <>
        <AuthContainer>
            <AuthCard>
                {!MainStore.getCheckAuthVariant() ? <AuthTitle onClick={()=>{MainStore.setCheckAuthVariant(true)}}>Регистрация</AuthTitle> : <AuthTitle onClick={()=>{MainStore.setCheckAuthVariant(false); console.log(MainStore.getCheckAuthVariant())}}>Вход</AuthTitle>}
                <TextField></TextField>
                <TextField></TextField>
                {MainStore.getCheckAuthVariant() ? <AuthButton>Зарегистрироваться</AuthButton> : <AuthButton>Войти</AuthButton>}
            </AuthCard>
        </AuthContainer>
    </>
}

export default AuthPage;