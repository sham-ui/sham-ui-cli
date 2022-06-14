import * as directives from 'sham-ui-directives';
import { createDI } from 'sham-ui';
import renderer from 'sham-ui-test-helpers';
//eslint-disable-next-line max-len
import RoutesCommonArticlesPage  from '../../../../../../src/components/routes/_common/articles/page.sfc';

it( 'renders correctly', async() => {
    const DI = createDI();

    DI.bind( 'title', {
        change() {}
    } );

    DI.bind( 'router', {
        storage: {
            params: {
                page: 2
            }
        },
        generate: () => '/'
    } );

    const meta = renderer( RoutesCommonArticlesPage, {
        DI,
        directives,
        loadData: () => Promise.resolve( {
            articles: [],
            meta: {
                total: 0,
                limit: 20
            }
        } )
    } );

    await new Promise( resolve => setImmediate( resolve ) );

    expect( meta.toJSON() ).toMatchSnapshot();
} );
