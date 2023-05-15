# WCH OtD API
This repository contains the code which is shared by all WCH OtD app netlify
functions. The exposed functions are in [`fetch.go`](/fetch.go).

## API Reference
The API can be queried at endpoints under `https://stories.workingclasshistory.com/api/{API version}/`

### `GET /api/v1/today`
A JSON-encoded list of events from our _On This Day in Working Class History_ series.

#### Query parameter
The `tz` query parameter can be specified, like `/api/v1/today?tz=America/New_York`. When the `tz` parameter is present, the results should be returned for the current day _in that time zone_, since it won't always be the same day in all time zones. **If `tz` is _not_ specified, results are returned for the `UTC` time zone**.
#### Possible Responses
##### `200 OK`
###### Headers
|key|value|
|---|-----|
|`Content-Type`|`application/json`|
|`Warning`|not present, `cache miss`, `cache failure`, or `cache set failure`. `cache set failure` may be present alongside `cache miss` or `cache failure`, separated by "`; `"|

###### Body
JSON-encoded **list** of objects like
|key|value type|example|
|---|----------|-------|
|`"title"`|string|`"Pullman strike begins"`|
|`"content"`|string, may contain HTML tags|`"On 11 May 1894, the Pullman railroad strike began in Chicago following the firing of three workers the previous day, called by Eugene Debs’ American Railroad Union (ARU).\nA month after it began, 400 ARU delegates from around the country met, and in defiance of Debs and their leadership agreed to boycott all Pullman railroad cars across the country in support of the workers in Chicago. The boycott began on June 26, when switchmen in Chicago refused to switch Pullman cars, and were fired. Their colleagues then walked out in their support.\nThe strike then spread down various railroads until soon all 26 roads out of Chicago were stopped, as were all of the transcontinental lines which carried Pullman cars. At its peak it was the biggest strike in US history to date, involving over 250,000 rail workers across 27 states and territories. That said, the union weakened its base of support by refusing to admit Black members, which enabled employers to hire some Black workers as strikebreakers. Despite this, some Black workers helped strikers blockade train tracks around Chicago.\nThen the US government intervened, granting an injunction against all strike activities across the country, and brought in federal troops. Thousands of US soldiers joined state militia and deputy marshals paid by the rail companies to attack the workers, shooting dozens. Still, the workers fought back, and workers around the country organised to call a general strike to force Pullman into arbitration. But these efforts were blocked by union leaders and eventually repression broke the strike.\nThis book tells its story, and that of other mass strikes in the US: https://shop.workingclasshistory.com/collections/books/products/strike-jeremy-brecher"`|
|`"more_info"`|string, may contain HTML tags|(usually empty)|
|`"excerpt"`|string, may contain HTML tags|First 250 visible characters from `"content"`: `"On 11 May 1894, the Pullman railroad strike began in Chicago following the firing of three workers the previous day, called by Eugene Debs’ American Railroad Union (ARU).\nA month after it began, 400 ARU delegates from around the country met, and in defiance of Debs and their leadership agreed to boyc..."`|
|`"author"`|string|`"Working Class History"`|
|`"url"`|string|The link to the full story.|
|`"media"`|`null` or `MediaInfo` object.|See [Media Info](#media-info)|

###### Media Info
If the `"media"` value is not `null`, it's an object like
|key|value type|example|
|---|----------|-------|
|`"url"`|string|`"https://workingclasshistory.com/wp-content/uploads/2022/12/2591134984_9b8023cf7f_k-1-1024x683.jpg"`|
|`"credit"`|string|`"Hossam el-Hamalawy (CC BY 2.0); https://www.flickr.com/photos/elhamalawy"`|
|`"caption"`|string|`"Workers' meeting at Mansoura-España garment factory."`|


##### `400 Bad Request`
Returned when the `tz` query parameter is specified, but not a valid time  zone.
In this case the response body is a JSON-encoded object with the single value `"error"`, like

```json
{"error": "invalid time zone"}
```

##### `500 Internal Server Error`
Returned if the data can't be retrieved from either the cache or the database due to some internal error. In this case the response body is always the following JSON-encoded message:
```json
{"error": "internal server error"}
```
#### Example Response
<details>
<summary>Long JSON response example</summary>

```json
[
  {
    "title": "Pullman strike begins",
    "content": "On 11 May 1894, the Pullman railroad strike began in Chicago following the firing of three workers the previous day, called by Eugene Debs’ American Railroad Union (ARU).\nA month after it began, 400 ARU delegates from around the country met, and in defiance of Debs and their leadership agreed to boycott all Pullman railroad cars across the country in support of the workers in Chicago. The boycott began on June 26, when switchmen in Chicago refused to switch Pullman cars, and were fired. Their colleagues then walked out in their support.\nThe strike then spread down various railroads until soon all 26 roads out of Chicago were stopped, as were all of the transcontinental lines which carried Pullman cars. At its peak it was the biggest strike in US history to date, involving over 250,000 rail workers across 27 states and territories. That said, the union weakened its base of support by refusing to admit Black members, which enabled employers to hire some Black workers as strikebreakers. Despite this, some Black workers helped strikers blockade train tracks around Chicago.\nThen the US government intervened, granting an injunction against all strike activities across the country, and brought in federal troops. Thousands of US soldiers joined state militia and deputy marshals paid by the rail companies to attack the workers, shooting dozens. Still, the workers fought back, and workers around the country organised to call a general strike to force Pullman into arbitration. But these efforts were blocked by union leaders and eventually repression broke the strike.\nThis book tells its story, and that of other mass strikes in the US: https://shop.workingclasshistory.com/collections/books/products/strike-jeremy-brecher",
    "more_info": "",
    "excerpt": "On 11 May 1894, the Pullman railroad strike began in Chicago following the firing of three workers the previous day, called by Eugene Debs’ American Railroad Union (ARU).\nA month after it began, 400 ARU delegates from around the country met, and in defiance of Debs and their leadership agreed to boyc...",
    "author": "Working Class History",
    "url": "https://stories.workingclasshistory.com/article/8336/pullman-strike-begins",
    "media": null
  },
  {
    "title": "Veracruz tenant prisoners freed",
    "content": "<p>On 11 May 1923, after months of agitation 150 mostly women rent strikers who had been jailed in the Mexican town of Veracruz the previous year were freed by the governor.</p> <p>The women had organised strikes in detention, and fought with prison guards, while workers outside threatened a general strike for their freedom.</p> <p>The tenants left the jail in groups of 10, the women wearing cream dresses and straw hats with red ribbons, while their supporters sang songs, shouted slogans and set off firecrackers. The group then paraded through the main streets of the city to the office of the renters' union, where they declared their commitment to continue their direct action against landlords.&nbsp;</p>",
    "more_info": "",
    "excerpt": "<p>On 11 May 1923, after months of agitation 150 mostly women rent strikers who had been jailed in the Mexican town of Veracruz the previous year were freed by the governor.</p> <p>The women had organised strikes in detention, and fought with prison guards, while workers outside threatened a general strike fo...</p>",
    "author": "Working Class History",
    "url": "https://stories.workingclasshistory.com/article/8337/veracruz-tenant-prisoners-freed",
    "media": null
  },
  {
    "title": "Police attack Latin Quarter",
    "content": "On 11 May 1968, French riot police began their assault  at 2:15 AM on Paris's Latin Quarter which had been occupied and barricaded by student protesters on the evening of May 10. Over the course of the night they eventually managed to evict the demonstrators, but the violence they employed against students and local residents provoked public anger and protests continued to grow.",
    "more_info": "",
    "excerpt": "On 11 May 1968, French riot police began their assault  at 2:15 AM on Paris's Latin Quarter which had been occupied and barricaded by student protesters on the evening of May 10. Over the course of the night they eventually managed to evict the demonstrators, but the violence they employed against stu...",
    "author": "Working Class History",
    "url": "https://stories.workingclasshistory.com/article/8338/police-attack-latin-quarter",
    "media": null
  },
  {
    "title": "Molaguero kidnapping",
    "content": "On 11 May 1972, a group called the Popular Revolutionary Organization 33 (OPR-33), the armed wing of the Uruguayan Anarchist Federation, kidnapped a shoe manufacturer named Sergio Molaguero whose workers were on strike. They received a ransom of $10 million, which they used to publicise internationally the crimes of the Uruguayan dictatorship. More info in this interesting pamphlet about the group: https://libcom.org/history/federacion-anarquista-uruguaya-fau-crisis-armed-struggle-dictatorship-1967-85",
    "more_info": "",
    "excerpt": "On 11 May 1972, a group called the Popular Revolutionary Organization 33 (OPR-33), the armed wing of the Uruguayan Anarchist Federation, kidnapped a shoe manufacturer named Sergio Molaguero whose workers were on strike. They received a ransom of $10 million, which they used to publicise interna...",
    "author": "Working Class History",
    "url": "https://stories.workingclasshistory.com/article/8339/molaguero-kidnapping",
    "media": null
  }
]
```
</details>

### `GET /api/v1/one_random_from_today`
Largely identical to `/api/v1/today` except that rather than returning all stories for today, it returns a single one selected at random.

#### Query parameter
The `tz` query parameter can be specified, like `/api/v1/today?tz=America/New_York`. When the `tz` parameter is present, the results should be returned for the current day _in that time zone_, since it won't always be the same day in all time zones. **If `tz` is _not_ specified, results are returned for the `UTC` time zone**.
#### Possible Responses
##### `200 OK`
###### Headers
|key|value|
|---|-----|
|`Content-Type`|`application/json`|
|`Warning`|not present, `cache miss`, `cache failure`, or `cache set failure`. `cache set failure` may be present alongside `cache miss` or `cache failure`, separated by "`; `"|

###### Body
A single JSON-encoded object (not a list), of the same format described in `/today`'s [Body documentation](#body).

##### `400 Bad Request`
Returned when the `tz` query parameter is specified, but not a valid time  zone.
In this case the response body is a JSON-encoded object with the single value `"error"`, like

```json
{"error": "invalid time zone"}
```

##### `500 Internal Server Error`
Returned if the data can't be retrieved from either the cache or the database due to some internal error. In this case the response body is always the following JSON-encoded message:
```json
{"error": "internal server error"}
```
#### Example Response
<details>
<summary>Long JSON response example</summary>

```json
{
  "title": "Ludlow Massacre",
  "content": "<p>On 20 April 1914, the Ludlow massacre took place when US troops opened fire with machine guns on a camp of striking miners and their families in Ludlow, Colorado.</p> <p>12,000 miners had gone out on strike the previous September against the Rockefeller family-owned Colorado Fuel and Iron Corporation (CF&amp;I) following the killing of an activist of the United Mine Workers of America (UMWA). They then demanded better safety at work, and to be paid in money, instead of company scrip (tokens which could only be redeemed in the company store).</p> <p>The Rockefellers evicted the striking miners and their families from their homes, and so they set up \"tent cities\" to live in collectively, which miners' wives helped run. Company thugs harassed strikers, and occasionally drove by camps riddling them with machine-gun fire, killing and injuring workers and their children.</p> <p>Eventually, the national guard was ordered to evict all the strike encampments, and the morning of April 20 they attacked the largest camp in Ludlow. They opened fire with machine guns on the tents of the workers and their families, who then returned fire. The main organiser of the camp, Louis Tikas, went to visit the officer in charge of the national guard to arrange a truce. But he was beaten to the ground then shot repeatedly in the back, killing him. That night, troops entered the camp and set fire to it, killing 11 children and two women, in addition to 13 other people who were killed in the fighting. The youngest victim was Elvira Valdez, aged just 3 months.</p> <p>Protests against the massacre broke out across the country, but the workers at CF&amp;I were defeated, and many of them were subsequently sacked and replaced with non-union miners. Over the course of the strike 66 people were killed, but no guardsmen or company thugs were prosecuted.</p>",
  "more_info": "",
  "excerpt": "<p>On 20 April 1914, the Ludlow massacre took place when US troops opened fire with machine guns on a camp of striking miners and their families in Ludlow, Colorado.</p> <p>12,000 miners had gone out on strike the previous September against the Rockefeller family-owned Colorado Fuel and Iron Corporation (CF...</p>",
  "author": "Working Class History",
  "url": "https://stories.workingclasshistory.com/article/9243/ludlow-massacre",
  "media":
  {
    "url": "https://workingclasshistory.com/wp-content/uploads/2023/02/04.20-x-Ludlow_striker_family_in_front_of_tent.jpg",
    "credit": "Denver Library/Wikimedia Commons",
    "caption": "Striking miners wives and children in the strikers' tent camp, 1914"
  }
}
```
</details>

## Configuration
This library relies on environment variables for configuration, which are set
in the netlify deployment settings, or in whatever environment the dependent
service is running in. Those environment variables are:

|key|default|
|---|-------|
|REACT_APP_BASEROW_TOKEN|(app crashes if not set; except when tests are running)|
|API_CACHE_LOCATION|localhost:6379|
|API_CACHE_PASSWORD||
|API_CACHE_USERNAME||
|API_CACHE_DB|0|

## Testing
Tests can be run locally by first starting a redis container in docker with
`docker run --publish 6379:6379 --detach redis:alpine` before running
`go test`.
