import { get, compact } from 'lodash';
import React from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import { reduxForm } from 'redux-form';
import editablePanel from './editablePanel';

import { no_op_action } from 'shared/utils';

import { PanelSwaggerField, PanelField } from 'shared/EditablePanel';

const EstimatesDisplay = props => {
  const fieldProps = {
    schema: props.PPMEstimateSchema,
    values: props.PPMEstimate,
  };

  return (
    <React.Fragment>
      <div className="editable-panel-column">
        <PanelSwaggerField
          title="DoD ID"
          fieldName="estimated_incentive"
          {...fieldProps}
        />
        <PanelSwaggerField fieldName="weight_estimate" {...fieldProps} />
        <PanelSwaggerField fieldName="planned_move_date" {...fieldProps} />
        <PanelSwaggerField fieldName="has_sit" {...fieldProps} />
        <PanelSwaggerField fieldName="days_in_storage" {...fieldProps} />
        max storage cost (two fields)
      </div>
      <div className="editable-panel-column">
        <PanelSwaggerField fieldName="pickup_postal_code" {...fieldProps} />
        <PanelSwaggerField
          fieldName="additional_pickup_postal_code"
          {...fieldProps}
        />
        <PanelSwaggerField
          fieldName="destination_postal_code"
          {...fieldProps}
        />
        <p>Distance estimate goes here</p>
      </div>
    </React.Fragment>
  );
};

const EstimatesEdit = props => {
  // const { schema } = props;
  return <React.Fragment>This is where the editing happens!</React.Fragment>;
};

const formName = 'ppm_estimate_and_details';

let PPMEstimatesPanel = editablePanel(EstimatesDisplay, EstimatesEdit);
PPMEstimatesPanel = reduxForm({ form: formName })(PPMEstimatesPanel);

function mapStateToProps(state) {
  return {
    // reduxForm
    formData: state.form[formName],
    initialValues: {},

    // Wrapper
    PPMEstimateSchema: get(
      state,
      'swagger.spec.definitions.PersonallyProcuredMovePayload',
    ),
    hasError: false,
    errorMessage: state.office.error,
    PPMEstimate: state.office.officePPMs[0], // unsure about this
    isUpdating: false,
  };
}

function mapDispatchToProps(dispatch) {
  return bindActionCreators(
    {
      update: no_op_action,
    },
    dispatch,
  );
}

export default connect(mapStateToProps, mapDispatchToProps)(PPMEstimatesPanel);
