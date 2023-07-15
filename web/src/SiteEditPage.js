import React from "react";
import {Button, Card, Col, Input, Row, Select} from "antd";
import * as SiteBackend from "./backend/SiteBackend";
import * as Setting from "./Setting";
import i18next from "i18next";

const {Option} = Select;

class SiteEditPage extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      classes: props,
      siteName: props.match.params.siteName,
      site: null,
    };
  }

  UNSAFE_componentWillMount() {
    this.getSite();
  }

  getSite() {
    SiteBackend.getSite(this.props.account.name, this.state.siteName)
      .then((site) => {
        this.setState({
          site: site,
        });
      });
  }

  parseSiteField(key, value) {
    if (["score"].includes(key)) {
      value = Setting.myParseInt(value);
    }
    return value;
  }

  updateSiteField(key, value) {
    value = this.parseSiteField(key, value);

    const site = this.state.site;
    site[key] = value;
    this.setState({
      site: site,
    });
  }

  renderSite() {
    return (
      <Card size="small" title={
        <div>
          {i18next.t("site:Edit Site")}&nbsp;&nbsp;&nbsp;&nbsp;
          <Button type="primary" onClick={this.submitSiteEdit.bind(this)}>{i18next.t("general:Save")}</Button>
        </div>
      } style={{marginLeft: "5px"}} type="inner">
        <Row style={{marginTop: "10px"}} >
          <Col style={{marginTop: "5px"}} span={2}>
            {i18next.t("general:Name")}:
          </Col>
          <Col span={22} >
            <Input value={this.state.site.name} onChange={e => {
              this.updateSiteField("name", e.target.value);
            }} />
          </Col>
        </Row>
        <Row style={{marginTop: "20px"}} >
          <Col style={{marginTop: "5px"}} span={2}>
            {i18next.t("general:Display name")}:
          </Col>
          <Col span={22} >
            <Input value={this.state.site.displayName} onChange={e => {
              this.updateSiteField("displayName", e.target.value);
            }} />
          </Col>
        </Row>
        <Row style={{marginTop: "20px"}} >
          <Col style={{marginTop: "5px"}} span={2}>
            {i18next.t("site:Domain")}:
          </Col>
          <Col span={22} >
            <Input value={this.state.site.domain} onChange={e => {
              this.updateSiteField("domain", e.target.value);
            }} />
          </Col>
        </Row>
        <Row style={{marginTop: "20px"}} >
          <Col style={{marginTop: "5px"}} span={2}>
            {i18next.t("site:Host")}:
          </Col>
          <Col span={22} >
            <Input value={this.state.site.host} onChange={e => {
              this.updateSiteField("host", e.target.value);
            }} />
          </Col>
        </Row>
        <Row style={{marginTop: "20px"}} >
          <Col style={{marginTop: "5px"}} span={2}>
            {i18next.t("site:SSL mode")}:
          </Col>
          <Col span={22} >
            <Select virtual={false} style={{width: "100%"}} value={this.state.site.sslMode} onChange={(value => {this.updateSiteField("sslMode", value);})}>
              {
                [
                  {id: "HTTP", name: "HTTP"},
                  {id: "HTTPS and HTTP", name: "HTTPS and HTTP"},
                  {id: "HTTPS Only", name: "HTTPS Only"},
                ].map((item, index) => <Option key={index} value={item.id}>{item.name}</Option>)
              }
            </Select>
          </Col>
        </Row>
        <Row style={{marginTop: "20px"}} >
          <Col style={{marginTop: "5px"}} span={2}>
            {i18next.t("site:SSL cert")}:
          </Col>
          <Col span={22} >
            <Input value={this.state.site.sslCert} onChange={e => {
              this.updateSiteField("sslCert", e.target.value);
            }} />
          </Col>
        </Row>
      </Card>
    );
  }

  submitSiteEdit() {
    const site = Setting.deepCopy(this.state.site);
    SiteBackend.updateSite(this.state.site.owner, this.state.siteName, site)
      .then((res) => {
        if (res) {
          Setting.showMessage("success", "Successfully saved");
          this.setState({
            siteName: this.state.site.name,
          });
          this.props.history.push(`/sites/${this.state.site.name}`);
        } else {
          Setting.showMessage("error", "failed to save: server side failure");
          this.updateSiteField("name", this.state.siteName);
        }
      })
      .catch(error => {
        Setting.showMessage("error", `failed to save: ${error}`);
      });
  }

  render() {
    return (
      <div>
        <Row style={{width: "100%"}}>
          <Col span={1}>
          </Col>
          <Col span={22}>
            {
              this.state.site !== null ? this.renderSite() : null
            }
          </Col>
          <Col span={1}>
          </Col>
        </Row>
        <Row style={{margin: 10}}>
          <Col span={2}>
          </Col>
          <Col span={18}>
            <Button type="primary" size="large" onClick={this.submitSiteEdit.bind(this)}>{i18next.t("general:Save")}</Button>
          </Col>
        </Row>
      </div>
    );
  }
}

export default SiteEditPage;
