import { createDI } from 'sham-ui';
import RoutesArticlePage  from '../../../../../src/components/routes/article/page.sfc';
import renderer from 'sham-ui-test-helpers';

it( 'renders correctly', () => {
    const DI = createDI();
    DI.bind( 'title', {
        change() {}
    } );
    DI.bind( 'router', {
        storage: {
            params: {
                page: 2
            },
            url: '',
            addWatcher() {}
        },
        generate: () => '/'
    } );
    DI.bind( 'store', {
        articleBySlug: jest.fn().mockReturnValueOnce(
            Promise.resolve( { meta: {} } )
        )
    } );
    const meta = renderer( RoutesArticlePage, {
        DI
    } );
    expect( meta.toJSON() ).toMatchSnapshot();
} );
