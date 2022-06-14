import RuntimePanel  from '../../../../../src/components/routes/server-info/runtime-panel.sfc';
import renderer from 'sham-ui-test-helpers';

it( 'renders correctly', () => {
    const meta = renderer( RuntimePanel, {} );
    expect( meta.toJSON() ).toMatchSnapshot();
} );
