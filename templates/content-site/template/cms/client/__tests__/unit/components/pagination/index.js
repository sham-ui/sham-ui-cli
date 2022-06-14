import * as directives from 'sham-ui-directives';
import PaginationIndex  from '../../../../src/components/pagination/index.sfc';
import renderer from 'sham-ui-test-helpers';

it( 'renders correctly', () => {
    const meta = renderer( PaginationIndex, {
        directives
    } );
    expect( meta.toJSON() ).toMatchSnapshot();
} );


it( 'pages', () => {
    const meta = renderer( PaginationIndex, {
        directives,
        total: 50,
        offset: 40,
        limit: 20
    } );
    expect( meta.toJSON() ).toMatchSnapshot();
} );
