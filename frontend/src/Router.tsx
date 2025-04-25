
import React from "react"
import { Route, Routes } from "react-router"
import AuthPage from "./components/Auth/AuthPage.tsx"
import { routes } from "./lists/routes.ts"
import '../src/style.css'
interface Props {}

const Router: React.FC<Props> = ()=>{
    return <>
        <Routes>
            <Route path={routes.auth} element={<AuthPage/>}></Route>
        </Routes>
    </>
}

export default Router;