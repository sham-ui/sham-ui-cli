import { createDI } from 'sham-ui';
import hrefto from 'sham-ui-router/lib/href-to';
import RoutesArticleNotFound  from '../../../../../src/components/routes/article/not-found.sfc';
import renderer from 'sham-ui-test-helpers';

it( 'renders correctly', () => {
    const DI = createDI();
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

    const meta = renderer( RoutesArticleNotFound, {
        DI,
        directives: {
            hrefto
        }
    } );
    expect( meta.toJSON() ).toMatchSnapshot();
} );
