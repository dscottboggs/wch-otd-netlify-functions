package wch_otd_netlify_functions

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/dscottboggs/attest"
)

const (
	testStr         = `<p>On 17 April 1975, a wave of anti-fascist anger erupted in Milan following the murder of socialist activist, Claudio Varalli, by a fascist the day before.</p> <p>In the morning, a demonstration of 50,000 workers and students marched through the city centre, before splitting into two separate marches which confronted both police and fascists. The Milan headquarters of the fascist Italian Social Movement (MSI) was attacked and had to be defended by police who fired tear gas and, at one point, drove a truck into a group of protesters, killing 27-year-old anti-fascist teacher, Giannino Zibecchi.</p> <p>Attacks by anti-fascists carried on throughout the day: local MSI politician, Cesare Biglia, and fascist trade unionist, Rodolfo Mersi,&nbsp; were both beaten up. Anti-fascists raided numerous fascist offices and bars known to be frequented by fascists were either burnt down or otherwise destroyed.</p> <p>Protests spread to other cities, but the most militant continued to be in Milan where dozens of demonstrations, big and small, continued for days. The offices of two fascists - Gastone Nencioni (a senator), and Benito Bollati (a lawyer) - were attacked with molotov cocktails while the office of far-right CISNAL union was also stormed. Attacks were also carried out on some offices related to social democratic parties.</p>"`
	exampleResponse = `{"id":9299,"order":"2025.00000000000000000000","title":"Operation Red Snake of the Paraná","year":"1975","month":"3","day":"20","time":"","description":"On 20 March 1975, Operation Red Snake of the Paraná began in Argentina, when the government of Isabel Perón sent hundreds of police and troops into the town of Villa Constitución to break the organisation of militant industrial workers. They arrested 307 workers, but the working class fought back, going on strike and occupying their plants until the detainees were released. The occupations lasted until March 26 when they were forcibly broken up by police. Over the next two months the government continued to arrest, blacklist and kill workers in the name of fighting \"subversion\". This sort of repression is well-known under the military dictatorship but much less so under Peron, who had the support of most unions.","social_description":null,"media":"","media_credit":"","media_caption":"","sources":"https://lapalabracaliente.wordpress.com/2015/03/19/40-anos-del-operativo-serpiente-colorada-del-parana-el-terrorismo-de-estado-antes-del-golpe/ accessed 13/03/2019","latitude":"","longitude":"","geotag_info":"","geotag_description":"","visitor_info":"","author_name":"Working Class History","author_url":"https://workingclasshistory.com","author_email":"","created_at":"2022-07-19","status":{"id":17945,"value":"published","color":"dark-gray"},"reviwer_id":[],"translations":[],"categories":[],"tags":[{"id":7664,"value":"military"},{"id":7686,"value":"Argentina"},{"id":7690,"value":"strikes"},{"id":7724,"value":"repression"},{"id":7770,"value":"killings"},{"id":7799,"value":"occupations"},{"id":7815,"value":"police"},{"id":11578,"value":"Isabel Perón"},{"id":11579,"value":"Partido Justicialista"},{"id":11580,"value":"Villa Constitución"}],"preview_text":null,"more_info":"","extra_media_media":[],"extra_media_caption":[],"extra_media_credit":[],"extra_media":[],"updated_at":"2022-07-19T22:01:05.113607Z","podcast_url":"","books_url":"","merch_url":"","spreadsheet_ref":"1960"}`
)

func TestExcerpt(t *testing.T) {
	test := attest.New(t)
	response := DbResponseRow{Description: testStr}
	result := test.EatError(response.Excerpt())
	test.Equals(result, `<p>On 17 April 1975, a wave of anti-fascist anger erupted in Milan following the murder of socialist activist, Claudio Varalli, by a fascist the day before.</p> <p>In the morning, a demonstration of 50,000 workers and students marched through the city centre, before splitting into two separate marches which confronted both police and fascists. The Milan headquarters of the fascist Italian Social Movement (MSI) was attacked and had to be defended by police who fired tear gas and, at one point, drove a truck into a group of protesters, killing 27-year-old anti-fascist teacher, Giannino Zibecchi.</p> <p>Att...</p>`)
	response = DbResponseRow{Description: `<p>invalid html</p></p>`}
	errResult, err := response.Excerpt()
	test.Equals("", errResult)
	test.NotNil(err, "invalid html error was nil")
	test.Compares(err, fmt.Errorf("unbalanced tags"))
}

func TestUrlEncodedTitle(t *testing.T) {
	var (
		test     = attest.New(t)
		response = DbResponseRow{Title: "Concha Liaño dies"}
	)
	test.Equals(response.UrlEncodedTitle(), "concha-lia%C3%B1o-dies")
}

func TestArticleUrl(t *testing.T) {
	test := attest.New(t)
	response := DbResponseRow{
		ID:    9299,
		Title: "Test Title",
	}
	articleUrl, err := response.ArticleUrl()
	test.Nil(err)
	test.Equals(articleUrl.String(), "https://stories.workingclasshistory.com/article/9299/test-title")
}

func TestSerializeAndTransform(t *testing.T) {
	var (
		example *DbResponseRow = new(DbResponseRow)
		test                   = attest.NewImmediate(t)
		err                    = json.NewDecoder(strings.NewReader(exampleResponse)).Decode(example)
	)
	test.StopIf(err)
	response := test.FailOnError(example.Transform()).(*OurResponse)
	test.Equals(response.Title, "Operation Red Snake of the Paraná")
	test.Equals(response.Content, "On 20 March 1975, Operation Red Snake of the Paraná began in Argentina, when the government of Isabel Perón sent hundreds of police and troops into the town of Villa Constitución to break the organisation of militant industrial workers. They arrested 307 workers, but the working class fought back, going on strike and occupying their plants until the detainees were released. The occupations lasted until March 26 when they were forcibly broken up by police. Over the next two months the government continued to arrest, blacklist and kill workers in the name of fighting \"subversion\". This sort of repression is well-known under the military dictatorship but much less so under Peron, who had the support of most unions.")
	test.Equals(response.MoreInfo, "")
	test.Equals(response.Excerpt, "On 20 March 1975, Operation Red Snake of the Paraná began in Argentina, when the government of Isabel Perón sent hundreds of police and troops into the town of Villa Constitución to break the organisation of militant industrial workers. They arrested 307 workers, but the working class fought back, going on strike and occupying their plants until the detainees were released. The occupations lasted until March 26 when they were forcibly broken up by police. Over the next two months the government continued to arrest, blacklist and kill workers in the name of fighting \"subversion\". This sort o...")
	test.NotNil(response.Author, "response Author was nil")
	test.Equals(response.Author.Name, "Working Class History")
	test.Equals(response.Author.Url, "https://workingclasshistory.com")
	test.Equals(response.Author.Email, "")
	test.Equals(response.Url, "https://stories.workingclasshistory.com/article/9299/operation-red-snake-of-the-paran%C3%A1")
	test.Attest(response.Media == nil, "media was not nil? %#+v", response.Media)
}
