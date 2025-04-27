interface IRoutes {
  auth: string;
  start: string;
  check: string;
  checkRoute: string;
  recs: string;
}
export const routes: IRoutes = {
  auth: '/auth',
  start: '/',
  checkRoute: '/info/:id',
  check: '/info',
  recs: '/recommendations',
};
