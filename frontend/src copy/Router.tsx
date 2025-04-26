import React from 'react';
import { Route, Routes } from 'react-router';
import AuthPage from './components/Auth/AuthPage';
import { routes } from './lists/routes';
import '../src/style.css';
import Start from './components/Start/Start';
import Analityc from './components/Analitycs/Analityc';
import Recommendations from './components/Recommendations/Recommendations';
interface Props {}

const Router: React.FC<Props> = () => {
  return (
    <>
      <Routes>
        <Route path={routes.auth} element={<AuthPage />}></Route>
        <Route path={routes.start} element={<Start />}></Route>
        <Route path={routes.checkRoute} element={<Analityc />}></Route>
        <Route path={routes.recs} element={<Recommendations />}></Route>
      </Routes>
    </>
  );
};

export default Router;
