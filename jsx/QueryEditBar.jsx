import React from 'react';
import * as action from './action.js';
import Popover from 'react-popover';


export default class QueryEditBar extends React.Component {
  constructor(props, context) {
    super(props, context);
    this.handleSaveChanges = this.handleSaveChanges.bind(this);
    this.handleToggleSQLPreview = this.handleToggleSQLPreview.bind(this);

    // 1) Is there a way to move discard changes to here without using action?
    // 2) maybe move generateQuery from output to here?

    this.state = {
      isOpen: false,
      popOverText: "",
    };
  }

  togglePopover() {
    this.setState({ isOpen: !this.state.isOpen });
  }

  handleSaveChanges() {
    console.log("handleSaveChanges ");

    // TODO: must support multiple queries for multiple rows changes
    var query = this.props.generateQuery();
    action.executeQuery(query);
  }

  handleToggleSQLPreview() {
    console.log("handleSQLPreview");
    if (this.state.isOpen){
      this.setState({isOpen: false});
    } else {
      var query = this.props.generateQuery();
      query = query.split(";").join("\n");

      this.setState({
        popOverText: query,
        isOpen: true,
      });
    }
  }

  render() {
    var positionStyle = { top: this.props.dragBarPosition + 'px' };

    return (
      <div id="query_edit_bar">
        <button className="save_changes" onClick={this.handleSaveChanges.bind(this)} style={positionStyle}>Save Changes</button>
        <button className="discard_changes" onClick={this.props.onHandleDiscardChanges} style={positionStyle}>Discard Changes</button>
        <div className="row_number" style={positionStyle}>{this.props.numberOfRowsEdited} edited rows</div>

        <Popover
          isOpen={this.state.isOpen}
          body={this.state.popOverText}
          preferPlace={"right"}
          target={"sql_preview"}
          targetElement={"sql_preview"}
          tipSize={10}>
          <div
            className="sql_preview"
            onClick={this.handleToggleSQLPreview.bind(this)}
            style={positionStyle}>
              {!this.state.isOpen ? "Show SQL Preview" :
                                    "Hide SQL Preview" }
          </div>
        </Popover>
      </div>
    );
  }
}
