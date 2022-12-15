import setupUnsafe from 'sham-ui-unsafe';
import { createDI } from 'sham-ui';
import hrefto from 'sham-ui-router/lib/href-to';
import formatLocaleDate from '../../../../../src/filters/format-locale-date';
import RoutesArticlePost  from '../../../../../src/components/routes/article/post.sfc';
import renderer from 'sham-ui-test-helpers';

it( 'renders correctly', () => {
    const DI = createDI();
    setupUnsafe( DI );

    DI.bind( 'router', {
        generate: () => '/'
    } );

    const meta = renderer( RoutesArticlePost, {}, {
        DI,
        directives: {
            hrefto
        },
        filters: {
            formatLocaleDate
        }
    } );
    expect( meta.toJSON() ).toMatchSnapshot();
} );
